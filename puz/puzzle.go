package puz

import (
	"fmt"
	"slices"
)

const file_magic string = "ACROSS&DOWN"
const default_version string = "1.4\x00"
const SOLID_SQUARE byte = '.'
const DIAGRAMLESS_SOLID_SQUARE byte = ':'
const EMPTY_STATE_SQUARE byte = '-'
const EMPTY_SOLUTION_SQUARE byte = ' '

type PuzzleType uint16

const (
	Normal      PuzzleType = 0x0001
	Diagramless PuzzleType = 0x0401
)

type Direction int

const (
	Across = iota
	Down
)

type extraSection int

const (
	rebus          extraSection = iota // GRBS
	rebusTable                         // RTBL
	timer                              // LTIM
	markup                             // GEXT
	userRebusTable                     // RUSR
)

var sectionMap = map[string]extraSection{
	"GRBS": rebus,
	"RTBL": rebusTable,
	"LTIM": timer,
	"GEXT": markup,
	"RUSR": userRebusTable,
}

var sectionStrMap = map[extraSection]string{
	rebus:          "GRBS",
	rebusTable:     "RTBL",
	timer:          "LTIM",
	markup:         "GEXT",
	userRebusTable: "RUSR",
}

func GetSectionFromString(s string) (extraSection, bool) {
	section, ok := sectionMap[s]
	return section, ok
}

func GetStrFromSection(s extraSection) (string, bool) {
	section, ok := sectionStrMap[s]
	return section, ok
}

type MarkupSquare byte

const (
	None                MarkupSquare = 0x00
	PreviouslyIncorrect MarkupSquare = 0x10
	CurrentlyIncorrect  MarkupSquare = 0x20
	ContentGiven        MarkupSquare = 0x40
	SquareCircled       MarkupSquare = 0x80
)

type Puzzle struct {
	Title         string
	Author        string
	Copyright     string
	Notes         string
	version       string
	Board         Board
	expectedClues uint16
	clues         Clues
	Extras        extraSections
	puzzleType    PuzzleType
	scramble      scrambleData
	UnusedData    unused
}

func NewPuzzle(width uint8, height uint8) *Puzzle {
	return &Puzzle{
		"",
		"",
		"",
		"",
		default_version,
		NewBoard(width, height),
		0,
		make([]Clue, 0),
		extraSections{
			make([]extraSection, 0),
			make([]RebusEntry, 0),
			nil,
			make([]RebusEntry, 0),
		},
		Normal,
		scrambleData{
			0, //TODO: check that this is right
			0,
		},
		unused{
			make([]byte, 2),
			make([]byte, 12),
			make([]byte, 0),
			make([]byte, 0),
		},
	}
}

func (p *Puzzle) SetVersion(version string) error {
	bytes := []byte(version)

	if len(version) != 3 || bytes[1] != '.' {
		return ErrInvalidVersionFormat
	}

	p.version = version + "\x00"

	return nil
}

func (p *Puzzle) Version() string {
	return p.version[:3]
}

func (p *Puzzle) Clues() Clues {
	return p.clues
}

func (p *Puzzle) SetClues(clues Clues) {
	p.clues = clues
	p.expectedClues = uint16(len(clues))
}

func (p *Puzzle) GetClueByPos(x int, y int, dir Direction) (*Clue, bool) {
	for _, clue := range p.clues {
		if clue.Direction == dir && clue.StartX == x && clue.StartY == y {
			return &clue, true
		}
	}

	return nil, false
}

func (p *Puzzle) GetClueByNum(num int, dir Direction) (*Clue, bool) {
	for _, clue := range p.clues {
		if clue.Direction == dir && clue.Num == num {
			return &clue, true
		}
	}

	return nil, false
}

func (p *Puzzle) AddClue(clue Clue, validateBoardPos bool) bool {
	if validateBoardPos {
		if clue.Direction == Across && !p.Board.StartsAcrossWord(clue.StartX, clue.StartY) {
			return false
		} else if clue.Direction == Down && !p.Board.StartsDownWord(clue.StartX, clue.StartY) {
			return false
		}
	}

	_, ok := p.GetClueByPos(clue.StartX, clue.StartY, clue.Direction)
	if ok {
		return false
	}

	p.expectedClues++
	p.clues = append(p.clues, clue)
	p.clues.Sort()
	return true
}

func (p *Puzzle) RemoveClueByPos(x int, y int, dir Direction) {
	p.clues = slices.DeleteFunc(p.clues, func(c Clue) bool {
		if c.StartX == x && c.StartY == y && c.Direction == dir {
			p.expectedClues--
			return true
		}

		return false
	})
}

func (p *Puzzle) RemoveClueByNum(num int, dir Direction) {
	p.clues = slices.DeleteFunc(p.clues, func(c Clue) bool {
		if c.Num == num && c.Direction == dir {
			p.expectedClues--
			return true
		}

		return false
	})
}

func (p *Puzzle) ExpectedClues() int {
	return int(p.expectedClues)
}

func (p *Puzzle) AddMarkupBoard() {
	p.Extras.extraSectionOrder = append(p.Extras.extraSectionOrder, markup)
}

func (p *Puzzle) Scrambled() bool {
	return p.scramble.scrambledTag != 0
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

type Board [][]Cell

func NewBoard(width uint8, height uint8) [][]Cell {
	board := make([][]Cell, height)

	for y := range height {
		board[y] = make([]Cell, width)

		for x := range width {
			board[y][x] = Cell{
				EMPTY_SOLUTION_SQUARE,
				EMPTY_STATE_SQUARE,
				0x00,
				0x00,
			}
		}
	}

	return board
}

func (b Board) Height() int {
	return len(b)
}

func (b Board) Width() int {
	return len(b[0])
}

func (b Board) inBounds(x int, y int) bool {
	if x >= 0 && x < b.Width() && y >= 0 && y < b.Height() {
		return true
	}

	return false
}

func (b Board) IsBlackSquare(x int, y int) bool {
	return b[y][x].Value == SOLID_SQUARE || b[y][x].Value == DIAGRAMLESS_SOLID_SQUARE
}

func (b Board) GetWord(x int, y int, dir Direction) (string, bool) {
	if !b.inBounds(x, y) {
		return "", false
	}

	if b.IsBlackSquare(x, y) {
		return "", false
	}

	word := ""

	xOffset := x
	yOffset := y

	for {
		word += string(b[yOffset][xOffset].Value)

		if dir == Across {
			xOffset += 1
		} else {
			yOffset += 1
		}

		if !b.inBounds(xOffset, yOffset) || b.IsBlackSquare(xOffset, yOffset) {
			break
		}
	}

	return word, true
}

func (b Board) StartsAcrossWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if x == 0 || b.IsBlackSquare(x-1, y) {
		if x+1 < b.Width() && !b.IsBlackSquare(x+1, y) {
			return true
		}
	}

	return false
}

func (b Board) StartsDownWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if y == 0 || b.IsBlackSquare(x, y-1) {
		if y+1 < b.Height() && !b.IsBlackSquare(x, y+1) {
			return true
		}
	}

	return false
}

func (b Board) GetWords() []Word {
	var words []Word

	width := b.Width()
	nextWordNum := 1
	for y := range b.Height() {
		for x := range width {
			if b.IsBlackSquare(x, y) {
				continue
			}

			startsAcrossWord := b.StartsAcrossWord(x, y)
			startsDownWord := b.StartsDownWord(x, y)

			if startsAcrossWord {
				word, ok := b.GetWord(x, y, Across)
				if ok {
					words = append(words, Word{
						word,
						nextWordNum,
						x,
						y,
						Across,
					})
				}
			}

			if startsDownWord {
				word, ok := b.GetWord(x, y, Down)
				if ok {
					words = append(words, Word{
						word,
						nextWordNum,
						x,
						y,
						Down,
					})
				}
			}

			if startsAcrossWord || startsDownWord {
				nextWordNum++
			}
		}
	}

	return words
}

type Word struct {
	Word      string
	Num       int
	StartX    int
	StartY    int
	Direction Direction
}

type Cell struct {
	Value    byte
	State    byte
	RebusKey byte
	Markup   byte
}

type Clues []Clue

func (c Clues) Sort() {
	slices.SortStableFunc(c, func(a Clue, b Clue) int {
		diff := a.StartX - b.StartX

		if diff == 0 {
			diff = a.StartY - b.StartY
		}

		if diff == 0 {
			diff = int(a.Direction) - int(b.Direction)
		}

		return diff
	})
}

type Clue struct {
	Clue      string
	Num       int
	StartX    int
	StartY    int
	Direction Direction
}

func NewClue(clue string, num int, x int, y int, dir Direction) Clue {
	return Clue{
		clue,
		num,
		x,
		y,
		dir,
	}
}

type extraSections struct {
	extraSectionOrder []extraSection
	RebusTable        []RebusEntry
	Timer             *TimerData
	UserRebusTable    []RebusEntry
}

type RebusEntry struct {
	Key   int
	Value string
}

type TimerData struct {
	SecondsPassed int
	Running       bool
}

type scrambleData struct {
	scrambledTag      uint16
	scrambledChecksum uint16
}

type unused struct {
	reserved1  []byte
	reserved2  []byte
	Preamble   []byte
	Postscript []byte
}
