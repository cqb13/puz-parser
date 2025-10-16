package puz

import (
	"encoding/binary"
	"fmt"
)

const file_magic = "ACROSS&DOWN"

func LoadPuz(bytes []byte) (*Puzzle, error) {
	var puzzle Puzzle

	err := parseHeader(bytes, &puzzle)
	if err != nil {
		return nil, fmt.Errorf("Failed to read header: %s", err)
	}

	return nil, nil
}

func parseHeader(bytes []byte, puzzle *Puzzle) error {
	if len(bytes) < 51 {
		return fmt.Errorf("Not enough data, expected header length of 51 bytes, found %d", len(bytes))
	}

	checksum := parseShort(bytes[:2])
	fileMagic := bytes[2:13]

	if string(fileMagic) != file_magic {
		return fmt.Errorf("Unexpected file magic, expected ACROSS&DOWN, found %s", string(fileMagic))
	}

	//TODO: validate checksums
	cibChecksum := parseShort(bytes[14:16])
	maskedLowChecksum := bytes[16:20]
	maskedHighChecksum := bytes[20:24]

	version := string(bytes[24:27])
	puzzle.Metadata.Version = version

	scrambledChecksum := parseShort(bytes[30:32])
	puzzle.Metadata.ScrambledChecksum = scrambledChecksum

	width := bytes[44]
	height := bytes[45]
	puzzle.Width = width
	puzzle.Height = height

	clueCount := parseShort(bytes[46:48])
	puzzle.NumClues = clueCount

	bitmask := parseShort(bytes[48:50])
	puzzle.Metadata.Bitmask = bitmask

	scrambledTag := parseShort(bytes[50:52])
	puzzle.Metadata.ScrambledTag = scrambledTag

	_, _, _, _ = checksum, cibChecksum, maskedLowChecksum, maskedHighChecksum
	return nil
}

func parseShort(bytes []byte) uint16 {
	return binary.LittleEndian.Uint16(bytes)
}
