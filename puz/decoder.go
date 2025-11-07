package puz

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Parse .puz file data from bytes and return Puzzle
func DecodePuz(bytes []byte) (*Puzzle, error) {
	var puzzle Puzzle

	reader := newPuzzleReader(bytes)

	fileMagicIndex := reader.Index([]byte(file_magic))
	if fileMagicIndex == -1 {
		return nil, ErrMissingFileMagic
	}
	preamble, err := reader.Read(fileMagicIndex - 2)
	puzzle.preamble = preamble

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
		if errors.Is(err, ErrUknownExtraSectionName) {
			break
		}

		// when an extra section is found, but cant be read
		if err != nil {
			return nil, err
		}
	}

	postscript := reader.ReadRemaining()
	puzzle.postscript = postscript

	// bytes[len(preamble):len(reader.bytes)-len(postscript)] to ensure only the actual data is checksummed
	computedChecksums := computeChecksums(bytes[len(preamble):len(reader.bytes)-len(postscript)], puzzle.size, puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.Clues, puzzle.Notes)

	if foundChecksums.cibChecksum != computedChecksums.cibChecksum {
		return nil, ErrCIBChecksumMismatch
	}

	if foundChecksums.checksum != computedChecksums.checksum {
		return nil, ErrGlobalChecksumMismatch
	}

	if foundChecksums.maskedLowChecksum != computedChecksums.maskedLowChecksum {
		return nil, ErrMaskedLowChecksumMismatch
	}

	if foundChecksums.maskedHighChecksum != computedChecksums.maskedHighChecksum {
		return nil, ErrMaskedHighChecksumMismatch
	}

	return &puzzle, nil
}

func parseHeader(reader *puzzleReader, puzzle *Puzzle) (*checksums, error) {
	if reader.Len() < 52 {
		return nil, ErrUnreadableData
	}

	checksum, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	fileMagic := reader.ReadStr()
	if string(fileMagic) != file_magic {
		return nil, ErrMissingFileMagic
	}

	cibChecksum, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}

	maskedLowChecksum, err := reader.Read(4)
	if err != nil {
		return nil, err
	}
	maskedHighChecksum, err := reader.Read(4)
	if err != nil {
		return nil, err
	}

	version, err := reader.Read(4)
	if err != nil {
		return nil, err
	}
	puzzle.metadata.Version = string(version)

	// not used in most files
	reserved1, err := reader.Read(2)
	if err != nil {
		return nil, err
	}
	puzzle.reserved1 = reserved1

	scrambledChecksum, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.metadata.scrambledChecksum = scrambledChecksum

	// not used in most files
	reserved2, err := reader.Read(12)
	if err != nil {
		return nil, err
	}
	puzzle.reserved2 = reserved2

	width, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	height, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	puzzle.width = width
	puzzle.height = height
	puzzle.size = int(width) * int(height)

	clueCount, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.numClues = clueCount

	bitmask, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.metadata.Bitmask = bitmask

	scrambledTag, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.metadata.ScrambledTag = scrambledTag

	foundChecksums := checksums{
		checksum,
		cibChecksum,
		[4]byte(maskedLowChecksum),
		[4]byte(maskedHighChecksum),
	}

	return &foundChecksums, nil
}

func parseSolutionAndState(reader *puzzleReader, puzzle *Puzzle) error {
	expectedLen := reader.offset + puzzle.size*2

	if expectedLen > reader.Len() {
		return ErrUnreadableData
	}

	solution, err := parseBoard(reader, int(puzzle.width), int(puzzle.height))
	if err != nil {
		return err
	}

	state, err := parseBoard(reader, int(puzzle.width), int(puzzle.height))
	if err != nil {
		return err
	}

	puzzle.Solution = solution
	puzzle.State = state

	return nil
}

func parseBoard(reader *puzzleReader, width int, height int) ([][]byte, error) {
	var board [][]byte

	for range height {
		bytes, err := reader.Read(width)
		if err != nil {
			return nil, err
		}

		board = append(board, bytes)
	}

	return board, nil
}

func parseStringsSection(reader *puzzleReader, puzzle *Puzzle) error {
	title := reader.ReadStr()
	puzzle.Title = title
	author := reader.ReadStr()
	puzzle.Author = author
	copyright := reader.ReadStr()
	puzzle.Copyright = copyright

	var clues []string

	for range puzzle.numClues {
		clue := reader.ReadStr()
		clues = append(clues, clue)
	}

	if len(clues) != int(puzzle.numClues) {
		return ErrClueCountMismatch
	}

	puzzle.Clues = clues

	notes := reader.ReadStr()
	puzzle.Notes = notes

	return nil
}

func parseExtraSection(reader *puzzleReader, puzzle *Puzzle) error {
	sectionName, err := reader.Peek(4)
	if err != nil {
		return ErrUknownExtraSectionName
	}

	section, ok := GetSectionFromString(string(sectionName))
	if !ok {
		return ErrUknownExtraSectionName
	}

	// just for shifting offset on valid section name
	reader.Read(4)

	// length does not include null terminator
	length, err := reader.ReadShort()
	if err != nil {
		return err
	}

	checksum, err := reader.ReadShort()
	if err != nil {
		return err
	}

	data, err := reader.Read(int(length))

	computedChecksum := checksumRegion(data, 0x00)

	if checksum != computedChecksum {
		return ErrExtraSectionChecksumMismatch
	}

	if slices.Contains(puzzle.extraSectionOrder, section) {
		return ErrDuplicateExtraSection
	}

	puzzle.extraSectionOrder = append(puzzle.extraSectionOrder, section)

	switch section {
	case Rebus:
		board, err := parseExtraSectionBoard(data, puzzle)
		if err != nil {
			return err
		}
		puzzle.ExtraSections.RebusBoard = board
	case RebusTable:
		tbl, err := parseExtraSectionRebusTbl(data)
		if err != nil {
			return err
		}
		puzzle.ExtraSections.RebusTable = tbl
	case Timer:
		timer, err := parseExtraTimerSection(data)
		if err != nil {
			return err
		}
		puzzle.ExtraSections.Timer = timer
	case Markup:
		board, err := parseExtraSectionBoard(data, puzzle)
		if err != nil {
			return err
		}
		puzzle.ExtraSections.MarkupBoard = board
	case UserRebusTable:
		tbl, err := parseExtraSectionRebusTbl(data)
		if err != nil {
			return err
		}
		puzzle.ExtraSections.UserRebusTable = tbl

	default:
		return ErrUknownExtraSectionName
	}

	// skip null terminator at the end of a section
	reader.ReadByte()

	return nil
}

func parseExtraSectionBoard(bytes []byte, puzzle *Puzzle) ([][]byte, error) {
	if len(bytes) != puzzle.size {
		return nil, ErrUnreadableData
	}

	var board [][]byte
	for i := 0; i < puzzle.size; i += int(puzzle.width) {
		end := i + int(puzzle.width)
		board = append(board, bytes[i:end])
	}

	return board, nil
}

func parseExtraSectionRebusTbl(bytes []byte) ([]RebusEntry, error) {
	// last byte is a ; and should be ignored for proper splitting
	str := string(bytes[:len(bytes)-1])

	parts := strings.Split(str, ";")

	var entries []RebusEntry

	for _, part := range parts {
		data := strings.Split(part, ":")
		if len(data) != 2 {
			return nil, ErrUnreadableData
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

func parseExtraTimerSection(bytes []byte) (*TimerData, error) {
	str := string(bytes)

	parts := strings.Split(str, ",")
	if len(parts) != 2 {
		return nil, ErrUnreadableData
	}

	if len(parts[1]) != 1 {
		return nil, ErrUnreadableData
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
