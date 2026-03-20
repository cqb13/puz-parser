package puz

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const headerSize = 52 // Header size in bytes

// DecodePuz parses .puz file data from bytes and returns Puzzle
func DecodePuz(data []byte) (*Puzzle, error) {
	var puzzle Puzzle

	reader := newPuzzleReader(data)

	fileMagicIndex := reader.index([]byte(fileMagic))
	if fileMagicIndex == -1 {
		return nil, MissingFileMagicError
	}
	preamble, err := reader.read(fileMagicIndex - 2)
	puzzle.UnusedData.Preamble = preamble

	foundChecksums, err := parseHeader(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse header: %w", err)
	}

	err = parseSolutionAndState(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse solution and state: %w", err)
	}

	err = parseStringsSection(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse strings section: %w", err)
	}

	for range 5 {
		err = parseExtraSection(&reader, &puzzle)
		if errors.Is(err, UnkownExtraSectionNameError) {
			break
		}

		// when an extra section is found, but cant be read
		if err != nil {
			return nil, err
		}
	}

	postscript := reader.readRemaining()
	puzzle.UnusedData.Postscript = postscript

	// to ensure only the actual data is checksummed
	computedChecksums := computeChecksums(data[len(preamble):len(reader.data)-len(postscript)], puzzle.Board.Width()*puzzle.Board.Height(), puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.clues, puzzle.Notes, puzzle.version)

	err = validateChecksums(foundChecksums, computedChecksums)
	if err != nil {
		return nil, err
	}

	return &puzzle, nil
}

type puzzleReader struct {
	data   []byte
	offset int
}

func newPuzzleReader(data []byte) puzzleReader {
	return puzzleReader{
		data,
		0,
	}
}

func (r *puzzleReader) canRead(amount int) bool {
	return r.offset+amount <= len(r.data)
}

func (r *puzzleReader) read(amount int) ([]byte, error) {
	if !r.canRead(amount) {
		return nil, OutOfBoundsReadError
	}

	start := r.offset
	r.offset += amount
	return r.data[start:r.offset], nil
}

func (r *puzzleReader) peek(amount int) ([]byte, error) {
	if !r.canRead(amount) {
		return nil, OutOfBoundsReadError
	}

	start := r.offset
	return r.data[start : r.offset+amount], nil
}

func (r *puzzleReader) readStr() string {
	var data []byte

	for i := r.offset; i < len(r.data) && r.data[i] != 0x00; i++ {
		data = append(data, r.data[i])
		r.offset++
	}

	r.offset++

	return string(data)
}

func (r *puzzleReader) index(target []byte) int {
	index := bytes.Index(r.data, target)

	return index
}

func (r *puzzleReader) readByte() (byte, error) {
	b, err := r.read(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *puzzleReader) len() int {
	return len(r.data)
}

func (r *puzzleReader) readShort() (uint16, error) {
	b, err := r.read(2)
	if err != nil {
		return 0, err
	}
	return parseShort(b), nil
}

func (r *puzzleReader) readRemaining() []byte {
	return r.data[r.offset:len(r.data)]
}

func parseShort(data []byte) uint16 {
	return binary.LittleEndian.Uint16(data)
}

func parseHeader(reader *puzzleReader, puzzle *Puzzle) (*checksums, error) {
	if reader.len() < headerSize {
		return nil, UnreadableDataError
	}

	checksum, err := reader.readShort()
	if err != nil {
		return nil, err
	}
	magic := reader.readStr()
	if string(magic) != fileMagic {
		return nil, MissingFileMagicError
	}

	cibChecksum, err := reader.readShort()
	if err != nil {
		return nil, err
	}

	maskedLowChecksum, err := reader.read(4)
	if err != nil {
		return nil, err
	}
	maskedHighChecksum, err := reader.read(4)
	if err != nil {
		return nil, err
	}

	version, err := reader.read(4)
	if err != nil {
		return nil, err
	}
	puzzle.version = string(version)

	// not used in most files
	reserved1, err := reader.read(2)
	if err != nil {
		return nil, err
	}
	puzzle.UnusedData.reserved1 = reserved1

	scrambledChecksum, err := reader.readShort()
	if err != nil {
		return nil, err
	}
	puzzle.scramble.scrambledChecksum = scrambledChecksum

	// not used in most files
	reserved2, err := reader.read(12)
	if err != nil {
		return nil, err
	}
	puzzle.UnusedData.reserved2 = reserved2

	width, err := reader.readByte()
	if err != nil {
		return nil, err
	}
	height, err := reader.readByte()
	if err != nil {
		return nil, err
	}

	puzzle.Board = NewBoard(width, height)

	clueCount, err := reader.readShort()
	if err != nil {
		return nil, err
	}
	puzzle.expectedClues = clueCount

	bitmask, err := reader.readShort()
	if err != nil {
		return nil, err
	}
	if bitmask == uint16(Diagramless) {
		puzzle.PuzzleType = Diagramless
	} else {
		puzzle.PuzzleType = Normal
	}

	scrambledTag, err := reader.readShort()
	if err != nil {
		return nil, err
	}
	puzzle.scramble.scrambledTag = scrambledTag

	foundChecksums := checksums{
		checksum,
		cibChecksum,
		[4]byte(maskedLowChecksum),
		[4]byte(maskedHighChecksum),
	}

	return &foundChecksums, nil
}

func parseSolutionAndState(reader *puzzleReader, puzzle *Puzzle) error {
	width := puzzle.Board.Width()
	height := puzzle.Board.Height()
	size := width * height

	expectedLen := reader.offset + size*2

	if expectedLen > reader.len() {
		return UnreadableDataError
	}

	solution, err := parseBoard(reader, width, height)
	if err != nil {
		return err
	}

	state, err := parseBoard(reader, width, height)
	if err != nil {
		return err
	}

	for y := range height {
		for x := range width {
			puzzle.Board[y][x].Answer = solution[y][x]
			puzzle.Board[y][x].Guess = state[y][x]
		}
	}

	return nil
}

func parseBoard(reader *puzzleReader, width int, height int) ([][]byte, error) {
	board := make([][]byte, height)

	for y := range height {
		row, err := reader.read(width)
		if err != nil {
			return nil, err
		}
		board[y] = row
	}

	return board, nil
}

func parseStringsSection(reader *puzzleReader, puzzle *Puzzle) error {
	title := reader.readStr()
	puzzle.Title = title
	author := reader.readStr()
	puzzle.Author = author
	copyright := reader.readStr()
	puzzle.Copyright = copyright

	var clues []string

	for range puzzle.expectedClues {
		clue := reader.readStr()
		clues = append(clues, clue)
	}

	if len(clues) != int(puzzle.expectedClues) {
		return &ClueCountMismatchError{
			int(puzzle.expectedClues),
			len(clues),
		}
	}

	puzzle.clues = make([]Clue, puzzle.expectedClues)

	notes := reader.readStr()
	puzzle.Notes = notes

	if puzzle.expectedClues == 0 {
		return nil
	}

	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	nextClueNum := 1
	nextClueIndex := 0
	assigned := false
	needsAcrossNum := false
	needsDownNum := false

	for y := range height {
		for x := range width {
			if puzzle.Board.IsSolidSquare(x, y) {
				continue
			}

			assigned = false

			needsAcrossNum = puzzle.Board.StartsAcrossWord(x, y)
			needsDownNum = puzzle.Board.StartsDownWord(x, y)

			if needsAcrossNum {
				puzzle.clues[nextClueIndex] = NewClue(clues[nextClueIndex], nextClueNum, x, y, Across)
				assigned = true
				nextClueIndex += 1
			}

			if needsDownNum {
				puzzle.clues[nextClueIndex] = NewClue(clues[nextClueIndex], nextClueNum, x, y, Down)
				assigned = true
				nextClueIndex += 1
			}

			if assigned {
				nextClueNum += 1
			}
		}
	}

	return nil
}

func parseExtraSection(reader *puzzleReader, puzzle *Puzzle) error {
	sectionName, err := reader.peek(4)
	if err != nil {
		return UnkownExtraSectionNameError
	}

	section, ok := getSectionFromString(string(sectionName))
	if !ok {
		return UnkownExtraSectionNameError
	}

	// because of peek, if we here there is data so no err check needed, just for shifting offset on valid section name
	reader.read(4)

	// length does not include null terminator
	length, err := reader.readShort()
	if err != nil {
		return err
	}

	checksum, err := reader.readShort()
	if err != nil {
		return err
	}

	data, err := reader.read(int(length))

	computedChecksum := checksumRegion(data, 0x00)

	if checksum != computedChecksum {
		return &ExtraSectionChecksumMismatchError{
			checksum,
			computedChecksum,
			section,
		}
	}

	if slices.Contains(puzzle.Extras.extraSectionOrder, section) {
		return &DuplicateExtraSectionError{
			section,
		}
	}

	puzzle.Extras.extraSectionOrder = append(puzzle.Extras.extraSectionOrder, section)

	switch section {
	case RebusSection:
		height := puzzle.Board.Height()
		width := puzzle.Board.Width()
		board, err := parseExtraSectionBoard(data, width, height)
		if err != nil {
			return err
		}

		for y := range height {
			for x := range width {
				puzzle.Board[y][x].RebusKey = board[y][x]
			}
		}
	case RebusTableSection:
		tbl, err := parseExtraSectionRebusTbl(data)
		if err != nil {
			return err
		}
		puzzle.Extras.RebusTable = tbl
	case TimerSection:
		timer, err := parseExtraTimerSection(data)
		if err != nil {
			return err
		}
		puzzle.Extras.Timer = *timer
	case MarkupBoardSection:
		height := puzzle.Board.Height()
		width := puzzle.Board.Width()
		board, err := parseExtraSectionBoard(data, width, height)
		if err != nil {
			return err
		}

		for y := range height {
			for x := range width {
				puzzle.Board[y][x].Markup = board[y][x]
			}
		}
	case UserRebusTableSection:
		tbl, err := parseExtraSectionRebusTbl(data)
		if err != nil {
			return err
		}
		puzzle.Extras.UserRebusTable = tbl

	default:
		return UnkownExtraSectionNameError
	}

	// skip null terminator at the end of a section
	reader.readByte()

	return nil
}

func parseExtraSectionBoard(data []byte, width int, height int) ([][]byte, error) {
	size := width * height
	if len(data) != size {
		return nil, UnreadableDataError
	}

	var board [][]byte
	for i := 0; i < size; i += width {
		end := i + width
		board = append(board, data[i:end])
	}

	return board, nil
}

func parseExtraSectionRebusTbl(data []byte) ([]RebusEntry, error) {
	// last byte is a ; and should be ignored for proper splitting
	str := string(data[:len(data)-1])

	parts := strings.Split(str, ";")

	var entries []RebusEntry

	for _, part := range parts {
		data := strings.Split(part, ":")
		if len(data) != 2 {
			return nil, UnreadableDataError
		}

		rawKey := data[0]
		key, err := strconv.Atoi(strings.Trim(rawKey, " "))
		if err != nil {
			return nil, err
		}
		value := data[1]

		entry := RebusEntry{
			key + 1, // key is 1 less than what is in the GRBS board
			value,
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func parseExtraTimerSection(data []byte) (*TimerData, error) {
	str := string(data)

	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return nil, UnreadableDataError
	}

	if len(parts[1]) != 1 {
		return nil, UnreadableDataError
	}

	running := true

	if parts[1] == "1" {
		running = false
	}

	seccondsPassed, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	return &TimerData{
		seccondsPassed,
		running,
	}, nil
}

func validateChecksums(found, computed *checksums) error {
	if found.cibChecksum != computed.cibChecksum {
		return &ChecksumMismatchError{
			int(found.cibChecksum),
			int(computed.cibChecksum),
			cibChecksum,
		}
	}

	if found.checksum != computed.checksum {
		return &ChecksumMismatchError{
			int(found.checksum),
			int(computed.checksum),
			globalChecksum,
		}
	}

	if found.maskedLowChecksum != computed.maskedLowChecksum {
		return &ChecksumMismatchError{
			int(binary.LittleEndian.Uint32(found.maskedLowChecksum[:])),
			int(binary.LittleEndian.Uint32(computed.maskedLowChecksum[:])),
			maskedLowChecksum,
		}
	}

	if found.maskedHighChecksum != computed.maskedHighChecksum {
		return &ChecksumMismatchError{
			int(binary.LittleEndian.Uint32(found.maskedHighChecksum[:])),
			int(binary.LittleEndian.Uint32(computed.maskedHighChecksum[:])),
			maskedHighChecksum,
		}
	}

	return nil
}
