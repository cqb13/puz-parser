package puz

const file_magic = "ACROSS&DOWN"

type Puzzle struct {
	Width     uint8
	Height    uint8
	NumClues  uint16
	Solution  [][]string
	State     [][]string
	Title     string
	Author    string
	Copyright string
	Clues     []string
	Notes     string
	Metadata  metadata
}

type metadata struct {
	Version           string
	Bitmask           uint16
	ScrambledTag      uint16
	ScrambledChecksum uint16
}

type checksums struct {
	checksum           uint16
	cibChecksum        uint16
	maskedLowChecksum  [4]byte
	maskedHighChecksum [4]byte
}
