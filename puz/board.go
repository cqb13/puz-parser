package puz

import (
	"strings"
)

// A Board represents a crossword grid, it is made up of cells which contain answers, player guesses, and markup information.
type Board [][]Cell

type Word struct {
	Word      string    // The word extracted from the grid
	Num       int       // The calculated number of the word based on grid layout
	StartX    int       // The x position of the first letter
	StartY    int       // The y position of the first letter
	Direction Direction // The direction of the word
}

// A Cell is a square in a crossword grid
type Cell struct {
	Answer   byte // The answer letter
	Guess    byte // The letter guessed by the player (used for saving game state)
	RebusKey byte // Indicates a connection to a value in the rebus table with the same key
	Markup   byte // Indicates applied markup
}

// NewBoard returns a Board of width x height.
//
// Each cells answer is EmptySolutionSquare, guess is EmptyStateSquare, no markup or rebus is applied.
func NewBoard(width uint8, height uint8) Board {
	board := make([][]Cell, height)

	for y := range height {
		board[y] = make([]Cell, width)

		for x := range width {
			board[y][x] = Cell{
				EmptySolutionSquare,
				EmptyStateSquare,
				0x00,
				0x00,
			}
		}
	}

	return board
}

// NewBoardFromArr returns a Board with the same dimensions as the provided 2D byte slice. Returns a BoardWidthMismatchError if the rows in byteBoard are not all the same length.
//
// A cells answer is set to the corresponding value to the byte board, the guess is EmptyStateSquare unless the answer is a SOLID_SQUARE or DiagramlessSolidSquare in which case the guess will match it, no markup or rebus is applied.
func NewBoardFromArr(byteBoard [][]byte) (Board, error) {
	board := make([][]Cell, len(byteBoard))

	prevWdith := len(byteBoard[0])
	for y, row := range byteBoard {
		board[y] = make([]Cell, len(row))
		if len(row) != prevWdith {
			return nil, BoardWidthMismatchError
		}
		for x, value := range row {
			cell := Cell{
				EmptySolutionSquare,
				EmptyStateSquare,
				0x00,
				0x00,
			}

			if value == SolidSquare || value == DiagramlessSolidSquare {
				cell.Guess = value
			}

			cell.Answer = value

			board[y][x] = cell
		}
	}

	return board, nil
}

// Height returns the number of rows in the board.
func (b Board) Height() int {
	return len(b)
}

// Width returns the number of columns in the board.
func (b Board) Width() int {
	if len(b) == 0 {
		return 0
	}

	return len(b[0])
}

// inBounds reports if an (x, y) is within a boards bounds.
// Valid coordinates satisfy 0 <= x < b.Width() and 0 <= y < b.Height().
func (b Board) inBounds(x int, y int) bool {
	if x >= 0 && x < b.Width() && y >= 0 && y < b.Height() {
		return true
	}

	return false
}

// IsSolidSquare reports if a cell at (x, y) is SOLID_SQUARE or DiagramlessSolidSquare.
func (b Board) IsSolidSquare(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	return b[y][x].Answer == SolidSquare || b[y][x].Answer == DiagramlessSolidSquare
}

// GetWord returns the series of letters starting at (x, y) in the given direction. Continues until the edge of the board or until a solid square is hit.
//
// The word is only valid if the bool is true.
func (b Board) GetWord(x int, y int, dir Direction) (string, bool) {
	if !b.inBounds(x, y) {
		return "", false
	}

	if b.IsSolidSquare(x, y) {
		return "", false
	}

	var word strings.Builder

	xOffset := x
	yOffset := y

	for {
		word.WriteString(string(b[yOffset][xOffset].Answer))

		if dir == Across {
			xOffset += 1
		} else {
			yOffset += 1
		}

		if !b.inBounds(xOffset, yOffset) || b.IsSolidSquare(xOffset, yOffset) {
			break
		}
	}

	return word.String(), true
}

// StartsAcrossWord reports if an across word starts at (x, y).
//
// A word must either have the board edge or a solid square to its left and have at least 2 letters.
func (b Board) StartsAcrossWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if b.IsSolidSquare(x, y) {
		return false
	}

	if x == 0 || b.IsSolidSquare(x-1, y) {
		if x+1 < b.Width() && !b.IsSolidSquare(x+1, y) {
			return true
		}
	}

	return false
}

// StartsAcrossWord reports if a down word starts at (x, y).
//
// A word must either have the board edge or a solid square above it and have at least 2 letters.
func (b Board) StartsDownWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if b.IsSolidSquare(x, y) {
		return false
	}

	if y == 0 || b.IsSolidSquare(x, y-1) {
		if y+1 < b.Height() && !b.IsSolidSquare(x, y+1) {
			return true
		}
	}

	return false
}

// GetWords returns a list of Words from the board.
func (b Board) GetWords() []Word {
	var words []Word

	width := b.Width()
	nextWordNum := 1
	for y := range b.Height() {
		for x := range width {
			if b.IsSolidSquare(x, y) {
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
