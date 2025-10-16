package puz

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
	Metadata  metadata
}

type metadata struct {
	Version           string
	Bitmask           uint16
	ScrambledTag      uint16
	ScrambledChecksum uint16
}
