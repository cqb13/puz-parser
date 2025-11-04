package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cqb13/puz-parser/puz"
)

func main() {
	fp, err := os.Open("tests/test-files/Crossword-PreAndPost.puz")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	puzzle, err := puz.DecodePuz(bytes)
	if err != nil {
		panic(err)
	}

	board := puz.NewBoardFromBytes(puzzle.Solution)

	for _, row := range board.Board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}

	for _, clue := range puzzle.Clues {
		fmt.Println(clue)
	}

	words := board.GetWords(puz.ACROSS)

	fmt.Println("Across:")
	for _, word := range words {
		fmt.Println(word)
	}

	words = board.GetWords(puz.DOWN)

	fmt.Println("\nDown:")
	for _, word := range words {
		fmt.Println(word)
	}

	testBoard := puz.NewBoard(5, 5)

	ok := testBoard.PlaceWord("tests", 0, 0, puz.ACROSS)
	if !ok {
		fmt.Println("cant place")
		return
	}

	for _, row := range testBoard.Board {
		for _, cell := range row {
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}
}
