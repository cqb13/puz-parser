package puz

import "fmt"

func EncodePuz(puzzle *Puzzle) ([]byte, error) {
	writer := newByteWriter()

	writer.WriteBytes(puzzle.preamble)

	err := encodeHeader(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode header: %s", err)
	}

	err = encodeSolutionAndState(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode solution and state: %s", err)
	}

	err = encodeStringsSection(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode strings section: %s", err)
	}

	writer.WriteBytes(puzzle.postscript)

	bodyBytes := writer.Bytes()[len(puzzle.preamble) : len(writer.Bytes())-len(puzzle.postscript)]
	computedChecksums := computeChecksums(bodyBytes, puzzle.Size, puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.Clues, puzzle.Notes)

	preambleOffset := len(puzzle.preamble)
	err = writer.OverwriteShort(preambleOffset+0, computedChecksums.checksum)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert checksum: %s", err)
	}

	err = writer.OverwriteShort(preambleOffset+14, computedChecksums.cibChecksum)
	if err != nil {
		return nil, fmt.Errorf("Failed to insert CIB checksum: %s", err)
	}

	err = writer.OverWrite(preambleOffset+16, computedChecksums.maskedLowChecksum[:])
	if err != nil {
		return nil, fmt.Errorf("Failed to insert Masked Low checksum: %s", err)
	}

	err = writer.OverWrite(preambleOffset+20, computedChecksums.maskedHighChecksum[:])
	if err != nil {
		return nil, fmt.Errorf("Failed to insert Masked High checksum: %s", err)
	}

	return writer.Bytes(), nil
}

func encodeHeader(puzzle *Puzzle, writer *byteWriter) error {
	if len(puzzle.reserved1) != 2 {
		return fmt.Errorf("Incorrect amount of bytes in first reserved section, expected 2, found %d", len(puzzle.reserved1))
	}

	if len(puzzle.reserved2) != 12 {
		return fmt.Errorf("Incorrect amount of bytes in second reserved section, expected 12, found %d", len(puzzle.reserved2))
	}

	// placeholder for file checksum, computed and inserted later
	writer.WritePlaceholder(2)

	writer.WriteString(file_magic)

	// placeholder for cib, maskedLow, and maskedHigh checksums, computed and inserted later
	writer.WritePlaceholder(10)
	writer.WriteBytes([]byte(puzzle.metadata.Version)) // not using write str because it already has the null terminator
	writer.WriteBytes(puzzle.reserved1)
	writer.WriteShort(puzzle.metadata.scrambledChecksum)
	writer.WriteBytes(puzzle.reserved2)
	writer.WriteByte(puzzle.Width)
	writer.WriteByte(puzzle.Height)
	writer.WriteShort(puzzle.NumClues)
	writer.WriteShort(puzzle.metadata.Bitmask)
	writer.WriteShort(puzzle.metadata.ScrambledTag)

	return nil
}

func encodeSolutionAndState(puzzle *Puzzle, writer *byteWriter) error {
	if len(puzzle.Solution) != int(puzzle.Height) {
		return fmt.Errorf("Height mismatch, expected solution height of %d, found %d", puzzle.Height, len(puzzle.Solution))
	}

	if len(puzzle.State) != int(puzzle.Height) {
		return fmt.Errorf("Height mismatch, expected state height of %d, found %d", puzzle.Height, len(puzzle.State))
	}

	for i, row := range puzzle.Solution {
		if len(row) != int(puzzle.Width) {
			return fmt.Errorf("Width mismatch, expected width of %d in solution row %d, found %d", puzzle.Width, i+1, len(row))
		}
		writer.WriteBytes(row)
	}

	for i, row := range puzzle.State {
		if len(row) != int(puzzle.Width) {
			return fmt.Errorf("Width mismatch, expected width of %d in state row %d, found %d", puzzle.Width, i+1, len(row))
		}
		writer.WriteBytes(row)
	}

	return nil
}

func encodeStringsSection(puzzle *Puzzle, writer *byteWriter) error {
	if len(puzzle.Clues) != int(puzzle.NumClues) {
		return fmt.Errorf("Expected %d clues, found %d", puzzle.NumClues, len(puzzle.Clues))
	}

	writer.WriteString(puzzle.Title)
	writer.WriteString(puzzle.Author)
	writer.WriteString(puzzle.Copyright)
	for _, clue := range puzzle.Clues {
		writer.WriteString(clue)
	}
	writer.WriteString(puzzle.Notes)

	return nil
}
