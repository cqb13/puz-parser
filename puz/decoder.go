package puz

import (
	"fmt"
	"strings"
)

func DecodePuz(bytes []byte) (*Puzzle, error) {
	var puzzle Puzzle

	reader := NewByteReader(bytes)

	foundChecksums, err := parseHeader(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse header: %s", err)
	}

	err = parseSolutionAndState(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse solution and state: %s", err)
	}

	err = parseStringsSection(&reader, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse strings section: %s", err)
	}

	//TODO: validate checksums

	expectedCibChecksum := checksumRegion(bytes[44:52])

	if foundChecksums.cibChecksum != expectedCibChecksum {
		return nil, fmt.Errorf("CIB Checksum mismatch, expected %d, found %d", expectedCibChecksum, foundChecksums.cibChecksum)
	}

	_ = foundChecksums

	return nil, nil
}

func parseHeader(reader *ByteReader, puzzle *Puzzle) (*checksums, error) {
	if reader.Len() < 51 {
		return nil, fmt.Errorf("Not enough data, expected header length of 51 bytes, found %d", reader.Len())
	}

	checksum, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	fileMagic := reader.ReadStr()
	if string(fileMagic) != file_magic {
		return nil, fmt.Errorf("Unexpected file magic, expected ACROSS&DOWN, found %s", string(fileMagic))
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

	version, err := reader.Read(3)
	if err != nil {
		return nil, err
	}
	puzzle.Metadata.Version = string(version)

	// skips reserved space, not used in most files
	reader.Step()
	reader.Step()

	scrambledChecksum, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.Metadata.ScrambledChecksum = scrambledChecksum

	// skips space, not sure why
	reader.SetOffset(44)
	width, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	height, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	puzzle.Width = width
	puzzle.Height = height

	clueCount, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.NumClues = clueCount

	bitmask, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.Metadata.Bitmask = bitmask

	scrambledTag, err := reader.ReadShort()
	if err != nil {
		return nil, err
	}
	puzzle.Metadata.ScrambledTag = scrambledTag

	foundChecksums := checksums{
		checksum,
		cibChecksum,
		[4]byte(maskedLowChecksum),
		[4]byte(maskedHighChecksum),
	}

	return &foundChecksums, nil
}

func parseSolutionAndState(reader *ByteReader, puzzle *Puzzle) error {
	expectedLen := reader.offset + int((puzzle.Width*puzzle.Height)*2)

	if expectedLen > reader.Len() {
		return fmt.Errorf("Not enough data, expected at least %d bytes, found %d", expectedLen, reader.Len())
	}

	solution, err := parseBoard(reader, int(puzzle.Width), int(puzzle.Height))
	if err != nil {
		return err
	}

	state, err := parseBoard(reader, int(puzzle.Width), int(puzzle.Height))
	if err != nil {
		return err
	}

	puzzle.Solution = solution
	puzzle.State = state

	return nil
}

func parseBoard(reader *ByteReader, width int, height int) ([][]string, error) {
	var board [][]string

	for range height {
		bytes, err := reader.Read(width)
		if err != nil {
			return nil, err
		}

		row := strings.Split(string(bytes), "")
		board = append(board, row)
	}

	return board, nil
}

func parseStringsSection(reader *ByteReader, puzzle *Puzzle) error {
	title := reader.ReadStr()
	puzzle.Title = title
	author := reader.ReadStr()
	puzzle.Author = author
	copyright := reader.ReadStr()
	puzzle.Copyright = copyright

	var clues []string

	for range puzzle.NumClues {
		clue := reader.ReadStr()
		clues = append(clues, clue)
	}

	if len(clues) != int(puzzle.NumClues) {
		return fmt.Errorf("Not enough clues, expected %d clues, found %d", puzzle.NumClues, len(clues))
	}

	puzzle.Clues = clues

	notes := reader.ReadStr()
	puzzle.Notes = notes

	return nil
}
