package puz

import (
	"fmt"
)

const file_magic string = "ACROSS&DOWN"
const default_version string = "1.4\x00"
const min_word_len = 2
const SOLID_SQUARE byte = '.'
const DIAGRAMLESS_SOLID_SQUARE byte = ':'
const EMPTY_STATE_SQUARE byte = '-'
const EMPTY_SOLUTION_SQUARE byte = ' '

type Direction int

const (
	ACROSS = iota
	DOWN
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
	Clues         Clues
	Extras        *extraSections
	bitmask       uint16
	scramble      *scrambleData
	unusedData    *unused
}

/*
TODO: method to get the positions of the start of every word and the directions a word can go at a pos
*/
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
	if !b.inBounds(x, y) {
		return false
	}

	return b[y][x].Value == SOLID_SQUARE || b[y][x].Value == DIAGRAMLESS_SOLID_SQUARE
}

func (b Board) CellNeedsAcrossNumber(x int, y int) bool {
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

// TODO: add a method to get the markups square type
// TODO: rename Value to solution or answer maybe
type Cell struct {
	Value    byte
	State    byte
	RebusKey byte
	Markup   byte
}

func (b Board) CellNeedsDownNumber(x int, y int) bool {
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

// TODO: add a sort method to clues to sort based on pos if same pos then dir
type Clues []Clue

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

// TODO: remove this and just make a func in encoder
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

// TODO: remove this and just make a func in encoder
func (t *TimerData) ToBytes() []byte {
	runningRep := 0

	if !t.Running {
		runningRep = 1
	}

	return fmt.Appendf(nil, "%d,%d", t.SecondsPassed, runningRep)
}

type scrambleData struct {
	scrambledTag      uint16
	scrambledChecksum uint16
}

type unused struct {
	reserved1  []byte
	reserved2  []byte
	preamble   []byte
	postscript []byte
}
