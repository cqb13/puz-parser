package puz

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// EncodePuz encodes the data in puzzle to bytes that can be saved as a .puz file
func EncodePuz(puzzle *Puzzle) ([]byte, error) {
	writer := newPuzzleWriter()

	writer.writeBytes(puzzle.UnusedData.Preamble)

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

	writer.writeBytes(puzzle.UnusedData.Postscript)

	bodyBytes := writer.bytes()[len(puzzle.UnusedData.Preamble) : len(writer.bytes())-len(puzzle.UnusedData.Postscript)]
	computedChecksums := computeChecksums(bodyBytes, puzzle.Board.Width()*puzzle.Board.Height(), puzzle.Title, puzzle.Author, puzzle.Copyright, puzzle.clues, puzzle.Notes, puzzle.version)

	preambleOffset := len(puzzle.UnusedData.Preamble)
	err = writer.overwriteShort(preambleOffset+0, computedChecksums.checksum)
	if err != nil {
		return nil, err
	}

	err = writer.overwriteShort(preambleOffset+14, computedChecksums.cibChecksum)
	if err != nil {
		return nil, err
	}

	err = writer.overwrite(preambleOffset+16, computedChecksums.maskedLowChecksum[:])
	if err != nil {
		return nil, err
	}

	err = writer.overwrite(preambleOffset+20, computedChecksums.maskedHighChecksum[:])
	if err != nil {
		return nil, err
	}

	return writer.bytes(), nil
}

type puzzleWriter struct {
	buffer bytes.Buffer
}

func newPuzzleWriter() *puzzleWriter {
	return &puzzleWriter{
		bytes.Buffer{},
	}
}

func (w *puzzleWriter) writeString(str string) {
	w.buffer.WriteString(str)
	w.buffer.WriteByte(0x00)
}

func (w *puzzleWriter) writePlaceholder(amount int) {
	b := make([]byte, amount)

	for i := range b {
		b[i] = 0x00
	}

	w.buffer.Write(b)
}

func (w *puzzleWriter) writeBytes(bytes []byte) {
	w.buffer.Write(bytes)
}

func (w *puzzleWriter) writeShort(short uint16) {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	w.buffer.Write(b)
}

func (w *puzzleWriter) writeByte(b byte) error {
	return w.buffer.WriteByte(b)
}

func (w *puzzleWriter) bytes() []byte {
	return w.buffer.Bytes()
}

func (w *puzzleWriter) overwrite(offset int, newBytes []byte) error {
	data := w.buffer.Bytes()

	if offset < 0 || offset > len(data) || offset+len(newBytes) > len(data) {
		return OutOfBoundsWriteError
	}

	copy(data[offset:], newBytes)

	return nil
}

func (w *puzzleWriter) overwriteShort(offset int, short uint16) error {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	err := w.overwrite(offset, b)
	if err != nil {
		return err
	}

	return nil
}

func encodeHeader(puzzle *Puzzle, writer *puzzleWriter) error {
	// placeholder for file checksum, computed and inserted later
	writer.writePlaceholder(2)

	writer.writeString(fileMagic)

	// placeholder for cib, maskedLow, and maskedHigh checksums, computed and inserted later
	writer.writePlaceholder(10)
	writer.writeBytes([]byte(puzzle.version)) // not using write str because it already has the null terminator
	writer.writeBytes(puzzle.UnusedData.reserved1)
	writer.writeShort(puzzle.scramble.scrambledChecksum)
	writer.writeBytes(puzzle.UnusedData.reserved2)
	writer.writeByte(byte(puzzle.Board.Width()))
	writer.writeByte(byte(puzzle.Board.Height()))
	writer.writeShort(uint16(len(puzzle.clues)))
	writer.writeShort(uint16(puzzle.PuzzleType))
	writer.writeShort(puzzle.scramble.scrambledTag)

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
			solution[(y*width)+x] = puzzle.Board[y][x].Answer
			state[(y*width)+x] = puzzle.Board[y][x].Guess
		}
	}

	writer.writeBytes(solution)
	writer.writeBytes(state)
}

func encodeStringsSection(puzzle *Puzzle, writer *puzzleWriter) error {
	if len(puzzle.clues) != int(puzzle.expectedClues) {
		return &ClueCountMismatchError{
			int(puzzle.expectedClues),
			len(puzzle.clues),
		}
	}

	writer.writeString(puzzle.Title)
	writer.writeString(puzzle.Author)
	writer.writeString(puzzle.Copyright)
	for _, clue := range puzzle.clues {
		writer.writeString(clue.Clue)
	}
	writer.writeString(puzzle.Notes)

	return nil
}

func encodeExtraSections(puzzle *Puzzle, writer *puzzleWriter) error {
	for _, section := range puzzle.Extras.extraSectionOrder {
		var data []byte

		switch section {
		case RebusSection, MarkupBoardSection:
			height := puzzle.Board.Height()
			width := puzzle.Board.Width()
			size := height * width

			board := make([]byte, size)

			for y := range height {
				for x := range width {
					var val byte

					if section == RebusSection {
						val = puzzle.Board[y][x].RebusKey
					} else {
						val = puzzle.Board[y][x].Markup
					}

					board[(y*width)+x] = val
				}
			}

			data = board
		case RebusTableSection:
			if puzzle.Extras.RebusTable == nil {
				return MissingExtraSectionError
			}

			for _, entry := range puzzle.Extras.RebusTable {
				padding := ""
				if entry.Key-1 < 10 {
					padding = " "
				}
				data = fmt.Appendf(data, "%s%d:%s;", padding, entry.Key-1, entry.Value)
			}
		case TimerSection:
			runningRep := 0

			if !puzzle.Extras.Timer.Running {
				runningRep = 1
			}

			data = fmt.Appendf(data, "%d,%d", puzzle.Extras.Timer.SecondsPassed, runningRep)
		case UserRebusTableSection:
			if puzzle.Extras.UserRebusTable == nil {
				return MissingExtraSectionError
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

		writer.writeBytes([]byte(section.String()))
		writer.writeShort(sectionLength)
		writer.writeShort(checksum)
		writer.writeBytes(data)
		writer.writeByte(0x00)
	}

	return nil
}
