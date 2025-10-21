package main

import (
	"dev/cqb13/puz-reader/puz"
	"fmt"
	"io"
	"os"
)

func displayHex(bytes []byte) {
	for i, b := range bytes {
		fmt.Printf("%02X ", b)
		if (i+1)%16 == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Println()
}

func main() {
	var path = "tests/test-files/"

	fp, err := os.Open(path + "crossword.puz")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	readBytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	displayHex(readBytes)

	puzzle, err := puz.DecodePuz(readBytes)
	if err != nil {
		panic(err)
	}

	encodedBytes, err := puz.EncodePuz(puzzle)
	if err != nil {
		panic(err)
	}
	puzzle.Display()

	fmt.Println()
	displayHex(encodedBytes)

	return
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

	puzzle, err := puz.DecodePuz(bytes)
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
