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

type ExtraSection int

const (
	RebusSection          ExtraSection = iota // GRBS
	RebusTableSection                         // RTBL
	TimerSection                              // LTIM
	MarkupBoardSection                        // GEXT
	UserRebusTableSection                     // RUSR
)

var sectionMap = map[string]ExtraSection{
	"GRBS": RebusSection,
	"RTBL": RebusTableSection,
	"LTIM": TimerSection,
	"GEXT": MarkupBoardSection,
	"RUSR": UserRebusTableSection,
}

var sectionStrMap = map[ExtraSection]string{
	RebusSection:          "GRBS",
	RebusTableSection:     "RTBL",
	TimerSection:          "LTIM",
	MarkupBoardSection:    "GEXT",
	UserRebusTableSection: "RUSR",
}

func getSectionFromString(s string) (ExtraSection, bool) {
	section, ok := sectionMap[s]
	return section, ok
}

func (s ExtraSection) String() string {
	return sectionStrMap[s]
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
	PuzzleType    PuzzleType
	scramble      scrambleData
	UnusedData    unused
}

// NewPuzzle creates a new puzzle with an empty board
func NewPuzzle(width uint8, height uint8) *Puzzle {
	return NewPuzzleFromBoard(NewBoard(width, height))
}

// NewPuzzleFromBoard creates a new puzzle with the given board
// If RebusKey or Markup values were changed in the board the corresponding extra sections should added with AddExtraSection()
func NewPuzzleFromBoard(board Board) *Puzzle {
	return &Puzzle{
		"",
		"",
		"",
		"",
		default_version,
		board,
		0,
		make([]Clue, 0),
		extraSections{
			make([]ExtraSection, 0),
			make([]RebusEntry, 0),
			TimerData{
				0,
				false,
			},
			make([]RebusEntry, 0),
		},
		Normal,
		scrambleData{
			0,
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

// SetVersion changes the version of the crossword.
// A properly formatted version is 2 digits separated by a period, 'X.X'.
// The default version for new crosswords is 1.4, other notable versions are 1.2 which means a puzzle will not include the notes section in checksums, along with 2.0 which allows for non ASCII characters to be included.
// Returns ErrInvalidVersionFormat if the version is not 3 characters long or the middle character is not a '.'.
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

// GetClueByPos searches for a clue with matching x, y (indices on the game board) coordinates and word direction
func (p *Puzzle) GetClueByPos(x int, y int, dir Direction) (*Clue, bool) {
	for _, clue := range p.clues {
		if clue.Direction == dir && clue.StartX == x && clue.StartY == y {
			return &clue, true
		}
	}

	return nil, false
}

// GetClueByNum searches for a clue with matching clue number and word direction
func (p *Puzzle) GetClueByNum(num int, dir Direction) (*Clue, bool) {
	for _, clue := range p.clues {
		if clue.Direction == dir && clue.Num == num {
			return &clue, true
		}
	}

	return nil, false
}

// AddClue takes in a clue and adds it to the clue list, then sorts the clues.
// If validateBoardPos is true, a check will be performed to ensure that the position of the clue is the start of a word on the board.
// Returns ok if the clue with the same position and direction does not already exist and if the clue passes validation checks.
// Clues are sorted by their position on the board, if there are two clues associated with a position, the Across clue will take priority.
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

// RemoveClueByPos removes a clue with matching x, y (indices on the game board) coordinates and word direction if it exists
func (p *Puzzle) RemoveClueByPos(x int, y int, dir Direction) {
	p.clues = slices.DeleteFunc(p.clues, func(c Clue) bool {
		if c.StartX == x && c.StartY == y && c.Direction == dir {
			p.expectedClues--
			return true
		}

		return false
	})
}

// RemoveClueByNum removes a clue with matching clue number and word direction if it exists
func (p *Puzzle) RemoveClueByNum(num int, dir Direction) {
	p.clues = slices.DeleteFunc(p.clues, func(c Clue) bool {
		if c.Num == num && c.Direction == dir {
			p.expectedClues--
			return true
		}

		return false
	})
}

// ExpectedClues returns the number of expected clues.
// This should always match the amount of clues in the clue list.
func (p *Puzzle) ExpectedClues() int {
	return int(p.expectedClues)
}

//TODO: add docs explaining that even if data is set in extra sections, if they are not explicitly added to the list the section will not be encoded into a puzzle

// AddExtraSection appends the given section to the list of included extra sections if it isn't already included.
func (p *Puzzle) AddExtraSection(section ExtraSection) bool {
	if p.HasExtraSection(section) {
		return false
	}

	p.Extras.extraSectionOrder = append(p.Extras.extraSectionOrder, section)

	return true
}

// AddExtraSection appends the given section to the list of included extra sections if it isn't already included.
func (p *Puzzle) RemoveExtraSection(section ExtraSection) bool {
	index := slices.Index(p.Extras.extraSectionOrder, section)

	if index == -1 {
		return false
	}

	p.Extras.extraSectionOrder = append(p.Extras.extraSectionOrder[:index], p.Extras.extraSectionOrder[index+1:]...)

	return true
}

// HasExtraSection returns true if the given section is in the list of extra sections
func (p *Puzzle) HasExtraSection(section ExtraSection) bool {
	return slices.Contains(p.Extras.extraSectionOrder, section)
}

/*
Sorts extra sections to comply with standard order

1. RebusSection           GRBS
2. RebusTableSection      RTBL
3. TimerSection           LTIM
4. MarkupBoardSection     GEXT
5. UserRebusTableSection  RUSR
*/
func (p *Puzzle) SortExtraSections() {
	slices.SortFunc(p.Extras.extraSectionOrder, func(a ExtraSection, b ExtraSection) int {
		return int(a) - int(b)
	})
}

func (p *Puzzle) Scrambled() bool {
	return p.scramble.scrambledTag != 0
}

// Unscramble attempts to unscramble the puzzle using the key.
// A valid key is made of 4 non zero digits.
// Unscrambling will fail if the board is already unscrambled, an invalid key is provided,
// the key is incorrect, the board has non ASCII letters (a-z / A-Z), or the board has less than 12 valid letters.
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

// Scramble attempts to scramble the puzzle using the key.
// A valid key is made of 4 non zero digits.
// Scrambling will fail if the board is already scrambled, an invalid key is provided,
// the board has non ASCII letters (a-z / A-Z), or the board has less than 12 valid letters.
func (p *Puzzle) Scramble(key int) error {
	if p.Scrambled() {
		return ErrPuzzleIsScrambled
	}

	err := scramble(p, key)
	if err != nil {
		return fmt.Errorf("Failed to scramble crossword: %w", err)
	}

	return nil
}

type Clues []Clue

// Sort will sort clues by their position in the board.
// Clues with a x, y closer to 0, 0 will be earlier.
// If two clues have the same position, the across clue will be first.
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
	extraSectionOrder []ExtraSection
	RebusTable        []RebusEntry
	Timer             TimerData
	UserRebusTable    []RebusEntry
}

// TODO: add methods to add rebus entries to board
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
