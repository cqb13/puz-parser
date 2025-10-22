package main

import (
	"bytes"
	"dev/cqb13/puz-parser/puz"
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

	entries, err := os.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()

		puzzle, raw, err := loadPuzzle(path + name)
		if err != nil {
			fmt.Printf("Failed on %s : %s\n", name, err)
			continue
		}

		encodedBytes, err := puz.EncodePuz(puzzle)

		if !bytes.Equal(raw, encodedBytes) {
			fmt.Printf("Not equal: %s\n", name)
		}
	}
}

func loadPuzzle(path string) (*puz.Puzzle, []byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		return nil, nil, err
	}

	puzzle, err := puz.DecodePuz(bytes)
	if err != nil {
		return nil, nil, err
	}

	return puzzle, bytes, nil
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
