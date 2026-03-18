package puz

import (
	"fmt"
	"slices"
)

// TODO: make name formatting consistent
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

// MarkupSquare represents formatting that can be added to a Cell
type MarkupSquare byte

const (
	None                MarkupSquare = 0x00
	PreviouslyIncorrect MarkupSquare = 0x10
	CurrentlyIncorrect  MarkupSquare = 0x20
	ContentGiven        MarkupSquare = 0x40
	SquareCircled       MarkupSquare = 0x80
)

type Puzzle struct {
	Title         string        // The title of the crossword
	Author        string        // The authors of the crossword
	Copyright     string        // The copyright information for the crossword
	Notes         string        // Additional notes for the crossword
	version       string        // The puz format version
	Board         Board         // The crossword grid, contains answers and game state and other formatting
	expectedClues uint16        // The expected number of clues
	clues         Clues         // The clues for the crossword
	Extras        extraSections // Optional extra sections, RebusTable, Timer, and UserRebusTable
	PuzzleType    PuzzleType    // The puzzle type, either Normal or Diagramless
	scramble      scrambleData  // Contains information about the puzzles scramble
	UnusedData    unused        // Contains unused bytes from the puz format along with additional data from before and after the puz data
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
		return InvalidVersionFormatError
	}

	p.version = version + "\x00"

	return nil
}

// Version returns the file version
func (p *Puzzle) Version() string {
	return p.version[:3]
}

// Clues returns the crossword clues
func (p *Puzzle) Clues() Clues {
	return p.clues
}

// SetClues overrides all of the crossword clues with the provided clues
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

// GetCluesByDirection retrieves all clues with a direction
func (p *Puzzle) GetCluesByDirection(dir Direction) Clues {
	var clues Clues

	for _, clue := range p.clues {
		if clue.Direction == dir {
			clues = append(clues, clue)
		}
	}

	return clues
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
		return PuzzleIsUnscrambledError
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
		return PuzzleIsScrambledError
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
	Clue      string    // The clue itself
	Num       int       // The number associating the clue with a word
	StartX    int       // The X position of the start of the word in the board
	StartY    int       // The Y position of the start of the word in the board
	Direction Direction // The direction of the word
}

// NewClue returns a Clue initialized with the given clue text, number, position (x, y), and direction.
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
	extraSectionOrder []ExtraSection // The order to write extra sections when encoding
	RebusTable        []RebusEntry   // The rebus table
	Timer             TimerData      // The state of the timer
	UserRebusTable    []RebusEntry   // The rebus table guessed by the player
}

// TODO: add methods to add rebus entries to board
type RebusEntry struct {
	Key   int    // Key links the entry to a cell on the board
	Value string // The value
}

type TimerData struct {
	SecondsPassed int  // The number of seconds passed in the game
	Running       bool // Whether or not the timer is running
}

type scrambleData struct {
	scrambledTag      uint16 // Indicates weather or not the crossword board is scrambled
	scrambledChecksum uint16 // The checksum for the unscrambled grid
}

type unused struct {
	reserved1  []byte // The first reserved section, found in the puz header
	reserved2  []byte // The second reserved section, found in the puz header
	Preamble   []byte // Any data from before the start of the puz header
	Postscript []byte // Any data from after the expected end of puz data
}
