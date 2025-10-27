package puz

import (
	"bytes"
	"encoding/binary"
)

type puzzleReader struct {
	bytes  []byte
	offset int
}

func newPuzzleReader(bytes []byte) puzzleReader {
	return puzzleReader{
		bytes,
		0,
	}
}

func (r *puzzleReader) CanRead(amount int) bool {
	return r.offset+amount <= len(r.bytes)
}

func (r *puzzleReader) Read(amount int) ([]byte, error) {
	if !r.CanRead(amount) {
		return nil, ErrOutOfBoundsRead
	}

	start := r.offset
	r.offset += amount
	return r.bytes[start:r.offset], nil
}

func (r *puzzleReader) Peek(amount int) ([]byte, error) {
	if !r.CanRead(amount) {
		return nil, ErrOutOfBoundsRead
	}

	start := r.offset
	return r.bytes[start : r.offset+amount], nil
}

func (r *puzzleReader) ReadStr() string {
	var bytes []byte

	for i := r.offset; i < len(r.bytes) && r.bytes[i] != 0x00; i++ {
		bytes = append(bytes, r.bytes[i])
		r.offset++
	}

	r.offset++

	return string(bytes)
}

func (r *puzzleReader) Index(target []byte) int {
	index := bytes.Index(r.bytes, target)
	if index == -1 {
		return -1
	}

	return index
}

func (r *puzzleReader) ReadByte() (byte, error) {
	b, err := r.Read(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *puzzleReader) Len() int {
	return len(r.bytes)
}

func (r *puzzleReader) ReadShort() (uint16, error) {
	b, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return parseShort(b), nil
}

func (r *puzzleReader) ReadRemaining() []byte {
	return r.bytes[r.offset:len(r.bytes)]
}

func parseShort(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}
