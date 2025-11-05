package puz

import (
	"fmt"
)

const file_magic string = "ACROSS&DOWN"
const min_word_len = 2
const BLACK_SQUARE byte = '.'
const EMPTY_STATE_SQUARE byte = '-'
const EMPTY_SOLUTION_SQUARE byte = ' '

type ExtraSection int

const (
	RebusBoard     ExtraSection = iota // GRBS
	RebusTable                         // RTBL
	Timer                              // LTIM
	MarkupBoard                        // GEXT
	UserRebusTable                     // RUSR
)

var sectionMap = map[string]ExtraSection{
	"GRBS": RebusBoard,
	"RTBL": RebusTable,
	"LTIM": Timer,
	"GEXT": MarkupBoard,
	"RUSR": UserRebusTable,
}

var sectionStrMap = map[ExtraSection]string{
	RebusBoard:     "GRBS",
	RebusTable:     "RTBL",
	Timer:          "LTIM",
	MarkupBoard:    "GEXT",
	UserRebusTable: "RUSR",
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

type Markup byte

const (
	None                Markup = 0x00
	PreviouslyIncorrect Markup = 0x10
	CurrentlyIncorrect  Markup = 0x20
	ContentGiven        Markup = 0x40
	SquareCircled       Markup = 0x80
)

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
	Solution          Board
	State             Board
	extraSectionOrder []ExtraSection
	ExtraSections     ExtraSections
	metadata          metadata
	reserved1         []byte
	reserved2         []byte
	preamble          []byte
	postscript        []byte
}

//TODO: add clue to word, take an x and a y, and a direction if on the board that is a word, and the clue in the proper location

func (p *Puzzle) Scrambled() bool {
	if p.metadata.ScrambledTag == 0 {
		return false
	}

	return true
}

func (p *Puzzle) Unscramble(key int) error {
	if !p.Scrambled() {
		return ErrPuzzleIsUnscrambled
	}

	err := unscramble(p, key)
	if err != nil {
		return fmt.Errorf("Failed to unscramble crossword: %w", err)
	}

	return nil
}

func (p *Puzzle) Scramble(key int) error {
	if p.Scrambled() {
		return ErrPuzzleIsScrambled
	}

	err := scramble(p, key)
	if err != nil {
		return fmt.Errorf("Failed to unscramble crossword: %w", err)
	}

	return nil
}

func (p *Puzzle) GetMetadata() metadata {
	return p.metadata
}

func (p *Puzzle) SetVersion(version string) error {
	if len(version) != 3 {
		return ErrInvalidVersionFormat
	}

	p.metadata.Version = version + "\x00"

	return nil
}

// ExtraSections holds optional data sections. Any  may be nil if not set.
// TODO: make markupboard a type with methods to get markup squares
type ExtraSections struct {
	RebusBoard     [][]byte
	RebusTable     []RebusEntry
	Timer          *TimerData
	MarkupBoard    [][]byte
	UserRebusTable []RebusEntry
}

type RebusEntry struct {
	Key   int
	Value string
}

// Returns key-1. The key 1 greater than what it is in binary so it matches the key in the Rebus board
func (r *RebusEntry) GetRealKey() int {
	return r.Key - 1
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
