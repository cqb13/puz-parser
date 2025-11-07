package puz

import (
	"fmt"
	"slices"
)

const file_magic string = "ACROSS&DOWN"
const min_word_len = 2
const SOLID_SQUARE byte = '.'
const EMPTY_STATE_SQUARE byte = '-'
const EMPTY_SOLUTION_SQUARE byte = ' '

type ExtraSection int

const (
	Rebus          ExtraSection = iota // GRBS
	RebusTable                         // RTBL
	Timer                              // LTIM
	Markup                             // GEXT
	UserRebusTable                     // RUSR
)

var sectionMap = map[string]ExtraSection{
	"GRBS": Rebus,
	"RTBL": RebusTable,
	"LTIM": Timer,
	"GEXT": Markup,
	"RUSR": UserRebusTable,
}

var sectionStrMap = map[ExtraSection]string{
	Rebus:          "GRBS",
	RebusTable:     "RTBL",
	Timer:          "LTIM",
	Markup:         "GEXT",
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

type MarkupSquare byte

const (
	None                MarkupSquare = 0x00
	PreviouslyIncorrect MarkupSquare = 0x10
	CurrentlyIncorrect  MarkupSquare = 0x20
	ContentGiven        MarkupSquare = 0x40
	SquareCircled       MarkupSquare = 0x80
)

type Puzzle struct {
	Title             string
	Author            string
	Copyright         string
	Notes             string
	width             uint8
	height            uint8
	size              int
	numClues          uint16
	Clues             []string
	Solution          Board
	State             Board // State from the player solving the puzzle
	extraSectionOrder []ExtraSection
	ExtraSections     ExtraSections
	metadata          metadata
	reserved1         []byte
	reserved2         []byte
	preamble          []byte
	postscript        []byte
}

func (p *Puzzle) GetWidth() int {
	return int(p.width)
}

func (p *Puzzle) GetHeight() int {
	return int(p.height)
}

func (p *Puzzle) GetSize() int {
	return p.size
}

func (p *Puzzle) GetPreamble() []byte {
	return p.preamble
}

func (p *Puzzle) GetPostscript() []byte {
	return p.preamble
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

// Resets state and syncs solid squares from the solution board to the state board, fails if boards are not the same size
func (p *Puzzle) SyncStateWithSolution() error {
	if len(p.Solution) != len(p.State) {
		return ErrBoardHeightMismatch
	}

	for y, row := range p.Solution {
		if len(row) != len(p.State[y]) {
			return ErrBoardWidthMismatch
		}
		for x, cell := range row {
			if cell == SOLID_SQUARE {
				p.State[y][x] = SOLID_SQUARE
			} else {
				p.State[y][x] = EMPTY_STATE_SQUARE
			}
		}
	}

	return nil
}

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

// returns ok if the puzzle does not already have a rebus board
func (p *Puzzle) AddRebusBoard() bool {
	if p.ExtraSections.RebusBoard != nil || slices.Contains(p.extraSectionOrder, Rebus) {
		return false
	}

	board := make([][]byte, p.height)

	for y := range p.height {
		board[y] = make([]byte, p.width)
	}

	// Follows the expected order
	p.extraSectionOrder = slices.Insert(p.extraSectionOrder, 0, Rebus)
	p.ExtraSections.RebusBoard = board

	return true
}

type RebusBoard [][]byte

func (r RebusBoard) GetKeys() []int {
	var keys []int

	for _, row := range r {
		for _, key := range row {
			if key == 0x00 {
				continue
			}

			if !slices.Contains(keys, int(key)) {
				keys = append(keys, int(key))
			}
		}
	}

	return keys
}

func (r RebusBoard) GetNextKey() int {
	max := 0

	for _, row := range r {
		for _, key := range row {
			if key == 0x00 {
				continue
			}

			if int(key) > max {
				max = int(key)
			}
		}
	}

	return max + 1
}

type MarkupBoard [][]byte

func (m MarkupBoard) GetMarkupSquare(x int, y int) (MarkupSquare, bool) {
	switch m[y][x] {
	case 0x00:
		return None, true
	case 0x10:
		return PreviouslyIncorrect, true
	case 0x20:
		return CurrentlyIncorrect, true
	case 0x40:
		return ContentGiven, true
	case 0x80:
		return SquareCircled, true
	}

	return 0x00, false
}

// ExtraSections holds optional data sections. Any  may be nil if not set.
type ExtraSections struct {
	RebusBoard     RebusBoard
	RebusTable     []RebusEntry
	Timer          *TimerData
	MarkupBoard    MarkupBoard
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
