package puz

import "fmt"

const file_magic = "ACROSS&DOWN"

type ExtraSection int

const (
	GRBS ExtraSection = iota // Rebus data
	RTBL                     // Rebus solution table
	LTIM                     // Timer
	GEXT                     // Cell style attributes
	RUSR                     // User rebus entries
)

type GEXTValue byte

const (
	PreviouslyIncorrect = 0x10
	CurrentlyIncorrect  = 0x20
	ContentGiven        = 0x40
	SquareCircled       = 0x80
)

type Puzzle struct {
	Title             string
	Author            string
	Copyright         string
	Notes             string
	Width             uint8
	Height            uint8
	Size              int
	NumClues          uint16
	Clues             []string
	Solution          [][]byte
	State             [][]byte
	extraSectionOrder []ExtraSection
	ExtraSections     ExtraSections
	metadata          metadata
	reserved1         []byte
	reserved2         []byte
	preamble          []byte
	postscript        []byte
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

// ExtraSections holds optional data sections. Any field may be nil if not set.
type ExtraSections struct {
	GRBS [][]byte
	RTBL map[int]string
	LTIM *TimerData
	GEXT [][]byte
	RUSR map[int]string
}

type TimerData struct {
	SecondsPassed int
	Running       bool
}

type metadata struct {
	Version           string
	Bitmask           uint16
	ScrambledTag      uint16
	scrambledChecksum uint16
}
