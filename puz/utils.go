package puz

import (
	"encoding/binary"
	"errors"
)

type ByteReader struct {
	bytes  []byte
	offset int
}

func NewByteReader(bytes []byte) ByteReader {
	return ByteReader{
		bytes,
		0,
	}
}

func (r *ByteReader) Read(amount int) ([]byte, error) {
	if r.offset+amount > len(r.bytes) {
		return nil, errors.New("Out of bounds")
	}

	start := r.offset
	r.offset += amount
	return r.bytes[start:r.offset], nil
}

func (r *ByteReader) ReadStr() string {
	var bytes []byte

	for i := r.offset; i < len(r.bytes) && r.bytes[i] != 0x00; i++ {
		bytes = append(bytes, r.bytes[i])
		r.offset++
	}

	r.offset++

	return string(bytes)
}

func (r *ByteReader) ReadByte() (byte, error) {
	b, err := r.Read(1)
	if err != nil {
		return 0, err
	}
	return b[0], nil
}

func (r *ByteReader) Len() int {
	return len(r.bytes)
}

func (r *ByteReader) ReadShort() (uint16, error) {
	b, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return parseShort(b), nil
}

func (r *ByteReader) Step() {
	r.offset++
}

func (r *ByteReader) SetOffset(offset int) error {
	if offset < 0 || offset > len(r.bytes) {
		return errors.New("Invalid offset")
	}

	r.offset = offset
	return nil
}

func parseShort(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}
