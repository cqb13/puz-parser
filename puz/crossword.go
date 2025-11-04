package puz

import (
	"fmt"
)

const file_magic string = "ACROSS&DOWN"
const BLACK_SQUARE byte = '.'
const min_word_len = 2

type ExtraSection int

const (
	GRBS ExtraSection = iota // Rebus data
	RTBL                     // Rebus solution table
	LTIM                     // Timer
	GEXT                     // Cell style attributes
	RUSR                     // User rebus entries
)

var sectionMap = map[string]ExtraSection{
	"GRBS": GRBS,
	"RTBL": RTBL,
	"LTIM": LTIM,
	"GEXT": GEXT,
	"RUSR": RUSR,
}

var sectionStrMap = map[ExtraSection]string{
	GRBS: "GRBS",
	RTBL: "RTBL",
	LTIM: "LTIM",
	GEXT: "GEXT",
	RUSR: "RUSR",
}

func GetSectionFromString(s string) (ExtraSection, bool) {
	section, ok := sectionMap[s]
	return section, ok
}

func GetStrFromSection(s ExtraSection) (string, bool) {
	section, ok := sectionStrMap[s]
	return section, ok
}

type Direction int

const (
	ACROSS = iota
	DOWN
)

const (
	None                = 0x00
	PreviouslyIncorrect = 0x10
	CurrentlyIncorrect  = 0x20
	ContentGiven        = 0x40
	SquareCircled       = 0x80
)

type Clue struct {
	Clue      string
	Direction Direction
	WordNum   int
	WordStart struct {
		X int
		Y int
	}
}

// TODO: when encoding clues they must be sorted first
type Clues [][]Clue

type Puzzle struct {
	Title             string
	Author            string
	Copyright         string
	Notes             string
	Width             uint8
	Height            uint8
	Size              int
	numClues          uint16
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
	RTBL []RebusEntry
	LTIM *TimerData
	GEXT [][]byte
	RUSR []RebusEntry
}

type RebusEntry struct {
	Key   int
	Value string
}

func (r *RebusEntry) ToBytes() []byte {
	// Keys are stored as + 1 in entries to match with the board, must convert back
	padding := ""
	if r.Key-1 < 10 {
		padding = " "
	}
	return fmt.Appendf(nil, "%s%d:%s;", padding, r.Key-1, r.Value)
}

type TimerData struct {
	SecondsPassed int
	Running       bool
}

func (t *TimerData) ToBytes() []byte {
	runningRep := 0

	if !t.Running {
		runningRep = 1
	}

	return fmt.Appendf(nil, "%d,%d", t.SecondsPassed, runningRep)
}

type metadata struct {
	Version           string
	Bitmask           uint16
	ScrambledTag      uint16
	scrambledChecksum uint16
}
