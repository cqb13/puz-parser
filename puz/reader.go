package puz

import (
	"encoding/binary"
	"fmt"
)

func LoadPuz(bytes []byte) (*Puzzle, error) {
	var puzzle Puzzle

	readHeader(bytes, &puzzle)

	return nil, nil
}

// TODO: ensure there are enough bytes for header
func readHeader(bytes []byte, puzzle *Puzzle) error {
	checksum := parseShort(bytes[:2])
	fileMagic := bytes[2:13]
	_ = checksum
	fmt.Println(string(fileMagic))
	return nil
}

func parseShort(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}
