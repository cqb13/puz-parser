package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cqb13/puz-parser/puz"
)

func main() {
	fp, _ := os.Open("./tests/test-files/Crossword.puz")
	defer fp.Close()

	bytes, _ := io.ReadAll(fp)

	p, _ := puz.DecodePuz(bytes)

	for _, clue := range p.Clues() {
		dir := "D"

		if clue.Direction == puz.Across {
			dir = "A"
		}

		fmt.Printf("%d.%s: %s\n", clue.Num, dir, clue.Clue)
	}

	for _, row := range p.Board {
		for _, cell := range row {
			fmt.Printf("%c ", cell.State)
		}

		fmt.Printf("\n")
	}
}
