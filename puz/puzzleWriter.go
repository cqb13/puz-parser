package puz

import (
	"bytes"
	"encoding/binary"
)

type puzzleWriter struct {
	buffer bytes.Buffer
}

func newPuzzleWriter() *puzzleWriter {
	return &puzzleWriter{
		bytes.Buffer{},
	}
}

func (w *puzzleWriter) WriteString(str string) {
	w.buffer.WriteString(str)
	w.buffer.WriteByte(0x00)
}

func (w *puzzleWriter) WritePlaceholder(amount int) {
	b := make([]byte, amount)

	for i := range b {
		b[i] = 0x00
	}

	w.buffer.Write(b)
}

func (w *puzzleWriter) WriteBytes(bytes []byte) {
	w.buffer.Write(bytes)
}

func (w *puzzleWriter) WriteShort(short uint16) {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	w.buffer.Write(b)
}

func (w *puzzleWriter) WriteByte(b byte) {
	w.buffer.WriteByte(b)
}

func (w *puzzleWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

func (w *puzzleWriter) OverWrite(offset int, newBytes []byte) error {
	data := w.buffer.Bytes()

	if offset < 0 || offset > len(data) || offset+len(newBytes) > len(data) {
		return ErrOutOfBoundsWrite
	}

	copy(data[offset:], newBytes)

	return nil
}

func (w *puzzleWriter) OverwriteShort(offset int, short uint16) error {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	err := w.OverWrite(offset, b)
	if err != nil {
		return err
	}

	return nil
}
