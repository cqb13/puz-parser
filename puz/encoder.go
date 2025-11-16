package puz

import "fmt"

func EncodePuz(puzzle *Puzzle) ([]byte, error) {
	writer := newPuzzleWriter()

	writer.WriteBytes(puzzle.preamble)

	err := encodeHeader(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode header: %w", err)
	}

	err = encodeSolutionAndState(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode solution and state: %w", err)
	}

	err = encodeStringsSection(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode strings section: %w", err)
	}

	err = encodeExtraSections(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode extra sections: %w", err)
	}

	writer.WriteBytes(puzzle.postscript)

	bodyBytes := writer.Bytes()[len(puzzle.preamble) : len(writer.Bytes())-len(puzzle.postscript)]
	computedChecksums := computeChecksums(bodyBytes, puzzle.size, puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.Clues, puzzle.Notes)

	preambleOffset := len(puzzle.preamble)
	err = writer.OverwriteShort(preambleOffset+0, computedChecksums.checksum)
	if err != nil {
		return nil, err
	}

	err = writer.OverwriteShort(preambleOffset+14, computedChecksums.cibChecksum)
	if err != nil {
		return nil, err
	}

	err = writer.OverWrite(preambleOffset+16, computedChecksums.maskedLowChecksum[:])
	if err != nil {
		return nil, err
	}

	err = writer.OverWrite(preambleOffset+20, computedChecksums.maskedHighChecksum[:])
	if err != nil {
		return nil, err
	}

	return writer.Bytes(), nil
}

func encodeHeader(puzzle *Puzzle, writer *puzzleWriter) error {
	// placeholder for file checksum, computed and inserted later
	writer.WritePlaceholder(2)

	writer.WriteString(file_magic)

	// placeholder for cib, maskedLow, and maskedHigh checksums, computed and inserted later
	writer.WritePlaceholder(10)
	writer.WriteBytes([]byte(puzzle.metadata.Version)) // not using write str because it already has the null terminator
	writer.WriteBytes(puzzle.reserved1)
	writer.WriteShort(puzzle.metadata.scrambledChecksum)
	writer.WriteBytes(puzzle.reserved2)
	writer.WriteByte(puzzle.width)
	writer.WriteByte(puzzle.height)
	writer.WriteShort(puzzle.numClues)
	writer.WriteShort(puzzle.metadata.Bitmask)
	writer.WriteShort(puzzle.metadata.ScrambledTag)

	return nil
}

func encodeSolutionAndState(puzzle *Puzzle, writer *puzzleWriter) error {
	//TODO: specify which board failed in err
	board, err := joinBoard(puzzle.Solution, int(puzzle.width), int(puzzle.height))
	if err != nil {
		return err
	}
	writer.WriteBytes(board)

	board, err = joinBoard(puzzle.State, int(puzzle.width), int(puzzle.height))
	if err != nil {
		return err
	}
	writer.WriteBytes(board)

	return nil
}

func encodeStringsSection(puzzle *Puzzle, writer *puzzleWriter) error {
	if len(puzzle.Clues) != int(puzzle.numClues) {
		return ErrClueCountMismatch
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

// TODO: specify which section is missing
func encodeExtraSections(puzzle *Puzzle, writer *puzzleWriter) error {
	for _, section := range puzzle.extraSectionOrder {
		name, ok := GetStrFromSection(section)
		if !ok {
			return ErrUknownExtraSectionName
		}

		var data []byte

		switch section {
		case rebus:
			if puzzle.ExtraSections.RebusBoard == nil {
				return ErrMissingExtraSection
			}

			board, err := joinBoard(puzzle.ExtraSections.RebusBoard, int(puzzle.width), int(puzzle.height))
			if err != nil {
				return err
			}
			data = board
		case rebusTable:
			if puzzle.ExtraSections.RebusTable == nil {
				return ErrMissingExtraSection
			}

			for _, entry := range puzzle.ExtraSections.RebusTable {
				data = append(data, entry.ToBytes()...)
			}
		case timer:
			if puzzle.ExtraSections.Timer == nil {
				return ErrMissingExtraSection
			}
			data = puzzle.ExtraSections.Timer.ToBytes()
		case markup:
			if puzzle.ExtraSections.MarkupBoard == nil {
				return ErrMissingExtraSection
			}
			board, err := joinBoard(puzzle.ExtraSections.MarkupBoard, int(puzzle.width), int(puzzle.height))
			if err != nil {
				return err
			}
			data = board
		case userRebusTable:
			if puzzle.ExtraSections.UserRebusTable == nil {
				return ErrMissingExtraSection
			}

			for _, entry := range puzzle.ExtraSections.UserRebusTable {
				data = append(data, entry.ToBytes()...)
			}
		}

		sectionLength := uint16(len(data))
		checksum := checksumRegion(data, 0x00)

		// name str should not have a null terminator
		writer.WriteBytes([]byte(name))
		writer.WriteShort(sectionLength)
		writer.WriteShort(checksum)
		writer.WriteBytes(data)
		writer.WriteByte(0x00)
	}

	return nil
}

func joinBoard(board [][]byte, width int, height int) ([]byte, error) {
	if len(board) != height {
		return nil, ErrBoardHeightMismatch
	}
	var data []byte

	for _, row := range board {
		if len(row) != width {
			return nil, ErrBoardWidthMismatch
		}

		data = append(data, row...)
	}

	return data, nil
}
