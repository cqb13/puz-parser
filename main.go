package main

import (
	"dev/cqb13/puz-reader/puz"
	"fmt"
	"io"
	"os"
)

func main() {
	var path = "tests/test-files/"

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		puzzle, err := loadPuzzle(path + name)
		if err != nil {
			fmt.Printf("Failed on %s : %s\n", name, err)
			continue
		}

		if puzzle.Scrambled() {
			fmt.Println("cracking ", name)
			key := crackPuzzle(puzzle)

			if key != -1 {
				fmt.Printf("Cracked %s : %d\n", name, key)
			} else {
				fmt.Println("Failed to crack ", name)
			}
		}
	}
}

func loadPuzzle(path string) (*puz.Puzzle, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	puzzle, err := puz.DecodePuz(bytes, false)
	if err != nil {
		return nil, err
	}

	return puzzle, nil
}

func crackPuzzle(puzzle *puz.Puzzle) int {
	key := 0

	for key < 10000 {
		err := puzzle.Unscramble(key)
		if err == nil {
			return key
		}

		key++
	}

	return -1
}
