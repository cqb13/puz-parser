package puz

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type byteReader struct {
	bytes  []byte
	offset int
}

func newByteReader(bytes []byte) byteReader {
	return byteReader{
		bytes,
		0,
	}
}

func (r *byteReader) Read(amount int) ([]byte, error) {
	if r.offset+amount > len(r.bytes) {
		return nil, errors.New("Out of bounds")
	}

	start := r.offset
	r.offset += amount
	return r.bytes[start:r.offset], nil
}

func (r *byteReader) ReadStr() string {
	var bytes []byte

	for i := r.offset; i < len(r.bytes) && r.bytes[i] != 0x00; i++ {
		bytes = append(bytes, r.bytes[i])
		r.offset++
	}

	r.offset++

	return string(bytes)
}

func (r *byteReader) ReadByte() (byte, error) {
	b, err := r.Read(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *byteReader) Len() int {
	return len(r.bytes)
}

func (r *byteReader) ReadShort() (uint16, error) {
	b, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return parseShort(b), nil
}

func (r *byteReader) ReadRemaining() []byte {
	return r.bytes[r.offset:len(r.bytes)]
}

func parseShort(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}

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

func (w *byteWriter) WriteByte(b byte) {
	w.buffer.WriteByte(b)
}

func (w *byteWriter) Bytes() []byte {
	return w.buffer.Bytes()
}

func (w *byteWriter) OverWrite(offset int, newBytes []byte) error {
	data := w.buffer.Bytes()

	if offset < 0 || offset > len(data) {
		return fmt.Errorf("Offset %d out of range [0, %d]", offset, len(data))
	}

	if offset+len(newBytes) > len(data) {
		return fmt.Errorf("Overwrite would exceed buffer length")
	}

	copy(data[offset:], newBytes)

	return nil
}

func (w *byteWriter) OverwriteShort(offset int, short uint16) error {
	b := make([]byte, 2)

	binary.LittleEndian.PutUint16(b, short)

	err := w.OverWrite(offset, b)
	if err != nil {
		return err
	}

	return nil
}
