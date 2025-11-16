package main

import (
	"github.com/cqb13/puz-parser/puz"
	"os"
)

func main() {
	board := puz.NewBoard(5, 5)
	board[0][0] = puz.SOLID_SQUARE

	puzzle := puz.NewPuzzle(board, "test crossword", "cqb13")

	puzzle.AddRebusBoard()
	puzzle.ExtraSections.RebusBoard[0][1] = 1
	puzzle.AddRebusEntry(1, "Test")

	puzzle.AddMarkupBoard()
	puzzle.MarkupSquare(1, 1, puz.SquareCircled)

	encoded, err := puz.EncodePuz(puzzle)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile("temp.puz", encoded, 0644)
	if err != nil {
		panic(err)
	}
}
