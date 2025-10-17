package main

import (
	"dev/cqb13/puz-reader/puz"
	"fmt"
	"io"
	"os"
)

func main() {
	var path = "tests/test-files/nyt_locked.puz"

	fp, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	puzzle, err := puz.DecodePuz(bytes, false)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Scrambled: %v\n", puzzle.Scrambled())
	puzzle.Display()

	if puzzle.Scrambled() {
		key := crackPuzzle(puzzle)
		fmt.Println(key)
	}
}

func crackPuzzle(puzzle *puz.Puzzle) int {
	key := 0

	for key < 10000 {
		if key%1000 == 0 {
			fmt.Println("Checking", key)
		}

		err := puzzle.Unscramble(key)
		if err == nil {
			return key
		}

		key++
	}

	return -1
}
