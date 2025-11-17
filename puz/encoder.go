package puz

import "fmt"

func EncodePuz(puzzle *Puzzle) ([]byte, error) {
	writer := newPuzzleWriter()

	writer.WriteBytes(puzzle.unusedData.preamble)

	err := encodeHeader(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode header: %w", err)
	}

	encodeSolutionAndState(puzzle, writer)

	err = encodeStringsSection(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode strings section: %w", err)
	}

	err = encodeExtraSections(puzzle, writer)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode extra sections: %w", err)
	}

	writer.WriteBytes(puzzle.unusedData.postscript)

	bodyBytes := writer.Bytes()[len(puzzle.unusedData.preamble) : len(writer.Bytes())-len(puzzle.unusedData.postscript)]
	computedChecksums := computeChecksums(bodyBytes, puzzle.Board.Width()*puzzle.Board.Height(), puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.Clues, puzzle.Notes)

	preambleOffset := len(puzzle.unusedData.preamble)
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
	writer.WriteBytes([]byte(puzzle.version)) // not using write str because it already has the null terminator
	writer.WriteBytes(puzzle.unusedData.reserved1)
	writer.WriteShort(puzzle.scramble.scrambledChecksum)
	writer.WriteBytes(puzzle.unusedData.reserved2)
	writer.WriteByte(byte(puzzle.Board.Width()))
	writer.WriteByte(byte(puzzle.Board.Height()))
	writer.WriteShort(uint16(len(puzzle.Clues)))
	writer.WriteShort(uint16(puzzle.puzzleType))
	writer.WriteShort(puzzle.scramble.scrambledTag)

	return nil
}

func encodeSolutionAndState(puzzle *Puzzle, writer *puzzleWriter) {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	size := height * width

	solution := make([]byte, size)
	state := make([]byte, size)

	for y := range height {
		for x := range width {
			solution[(y*width)+x] = puzzle.Board[y][x].Value
			state[(y*width)+x] = puzzle.Board[y][x].State
		}
	}

	writer.WriteBytes(solution)
	writer.WriteBytes(state)
}

func encodeStringsSection(puzzle *Puzzle, writer *puzzleWriter) error {
	if len(puzzle.Clues) != int(puzzle.expectedClues) {
		return ErrClueCountMismatch
	}

	writer.WriteString(puzzle.Title)
	writer.WriteString(puzzle.Author)
	writer.WriteString(puzzle.Copyright)
	for _, clue := range puzzle.Clues {
		writer.WriteString(clue.Clue)
	}
	writer.WriteString(puzzle.Notes)

	return nil
}

// TODO: specify which section is missing
func encodeExtraSections(puzzle *Puzzle, writer *puzzleWriter) error {
	for _, section := range puzzle.Extras.extraSectionOrder {
		name, ok := GetStrFromSection(section)
		if !ok {
			return ErrUknownExtraSectionName
		}

		var data []byte

		switch section {
		case rebus, markup:

			height := puzzle.Board.Height()
			width := puzzle.Board.Width()
			size := height * width

			board := make([]byte, size)

			for y := range height {
				for x := range width {
					var val byte

					if section == rebus {
						val = puzzle.Board[y][x].RebusKey
					} else {
						val = puzzle.Board[y][x].Markup
					}

					board[(y*width)+x] = val
				}
			}

			data = board
		case rebusTable:
			if puzzle.Extras.RebusTable == nil {
				return ErrMissingExtraSection
			}

			for _, entry := range puzzle.Extras.RebusTable {
				padding := ""
				if entry.Key-1 < 10 {
					padding = " "
				}
				data = fmt.Appendf(data, "%s%d:%s;", padding, entry.Key-1, entry.Value)
			}
		case timer:
			if puzzle.Extras.Timer == nil {
				return ErrMissingExtraSection
			}

			runningRep := 0

			if !puzzle.Extras.Timer.Running {
				runningRep = 1
			}

			data = fmt.Appendf(data, "%d,%d", puzzle.Extras.Timer.SecondsPassed, runningRep)
		case userRebusTable:
			if puzzle.Extras.UserRebusTable == nil {
				return ErrMissingExtraSection
			}

			for _, entry := range puzzle.Extras.UserRebusTable {
				padding := ""
				if entry.Key-1 < 10 {
					padding = " "
				}
				data = fmt.Appendf(data, "%s%d:%s;", padding, entry.Key-1, entry.Value)
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
