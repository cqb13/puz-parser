package puz

import (
	"bytes"
	"encoding/binary"
)

type byteWriter struct {
	buffer bytes.Buffer
}

func newByteWriter() *byteWriter {
	return &byteWriter{
		bytes.Buffer{},
	}
}

func (w *byteWriter) WriteString(str string) {
	w.buffer.WriteString(str)
	w.buffer.WriteByte(0x00)
}

func (w *byteWriter) WritePlaceholder(amount int) {
	b := make([]byte, amount)

	for i := range b {
		b[i] = 0x00
	}

	w.buffer.Write(b)
}

func (w *byteWriter) WriteBytes(bytes []byte) {
	w.buffer.Write(bytes)
}

func (w *byteWriter) WriteShort(short uint16) {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	w.buffer.Write(b)
}

func (w *byteWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

func EncodePuz(puzzle *Puzzle) ([]byte, error) {
	writer := newByteWriter()

	writeHeader(puzzle, writer)

	return writer.Bytes(), nil
}

// TODO: validate reserved lengths
func writeHeader(puzzle *Puzzle, writer *byteWriter) {
	// placeholder for file checksum, computed and inserted later
	writer.WritePlaceholder(2)

	writer.WriteString(file_magic)

	// placeholder for cib, maskedLow, and maskedHigh checksums, computed and inserted later
	writer.WritePlaceholder(10)
	writer.WriteBytes(puzzle.reserved1)
	writer.WriteShort(puzzle.metadata.ScrambledChecksum)
	writer.WriteBytes(puzzle.reserved2)
}
