package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cqb13/puz-parser/puz"
)

func main() {
	fp, _ := os.Open("./tests/test-files/Crossword.puz")

	bytes, _ := io.ReadAll(fp)

	p, _ := puz.DecodePuz(bytes)

	for _, row := range p.Board {
		for _, cell := range row {
			fmt.Printf("%c ", cell.Value)
		}

		fmt.Printf("\n")
	}
}
