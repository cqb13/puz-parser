package main

import (
	"dev/cqb13/puz-reader/puz"
	"io"
	"os"
)

func main() {
	fp, err := os.Open("testing.puz")
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

	puzzle.Display()
}
