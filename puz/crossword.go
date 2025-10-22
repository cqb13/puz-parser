package puz

import "fmt"

const file_magic = "ACROSS&DOWN"

type Puzzle struct {
	Title      string
	Author     string
	Copyright  string
	Notes      string
	Width      uint8
	Height     uint8
	Size       int
	NumClues   uint16
	Clues      []string
	Solution   [][]byte
	State      [][]byte
	metadata   metadata
	reserved1  []byte
	reserved2  []byte
	preamble   []byte
	postscript []byte
}

func (p *Puzzle) Display() {
	fmt.Println(p.String())
}

func (p *Puzzle) String() string {
	str := fmt.Sprintf("Title: %s\nAuthor: %s\nCopyright: %s\nNotes: %s\nVersion: %s\nSize: %dx%d\nClues:\n", p.Title, p.Author, p.Copyright, p.Notes, p.metadata.Version, p.Width, p.Height)
	for i, clue := range p.Clues {
		str += fmt.Sprintf("\t%d. %s\n", i+1, clue)
	}
	str += "Solution:\n"
	for _, row := range p.Solution {
		str += "\t"
		for _, letter := range row {
			str += string(letter) + " "
		}
		str += "\n"
	}
	str += "State:\n"
	for _, row := range p.State {
		str += "\t"
		for _, letter := range row {
			str += string(letter) + " "
		}
		str += "\n"
	}
	return str
}

func (p *Puzzle) Scrambled() bool {
	if p.metadata.ScrambledTag == 0 {
		return false
	}

	return true
}

func (p *Puzzle) Unscramble(key int) error {
	if !p.Scrambled() {
		return fmt.Errorf("Puzzle is already unscrambled")
	}

	err := unscramble(p, key)
	if err != nil {
		return fmt.Errorf("Failed to unscramble crossword: %s", err)
	}

	return nil
}

func (p *Puzzle) Scramble(key int) error {
	if p.Scrambled() {
		return fmt.Errorf("Puzzle is already scrambled")
	}

	err := scramble(p, key)
	if err != nil {
		return fmt.Errorf("Failed to unscramble crossword: %s", err)
	}

	return nil
}

func (p *Puzzle) GetMetadata() metadata {
	return p.metadata
}

func (p *Puzzle) SetVersion(version string) error {
	if len(version) != 3 {
		return fmt.Errorf("Invalid version format, must be X.X")
	}

	p.metadata.Version = version + "\x00"

	return nil
}

type metadata struct {
	Version           string
	Bitmask           uint16
	ScrambledTag      uint16
	ScrambledChecksum uint16
}
