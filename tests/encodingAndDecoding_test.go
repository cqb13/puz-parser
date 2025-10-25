package tests

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

const test_files = "./test-files/"

func TestCorrectDecodeAndEncode(t *testing.T) {
	// normal 5x5 puzzle
	puz1 := "Crossword.puz"
	correctDecodeAndEncode(puz1, t)

	// normal 5x5 puzzle with preamble and postscript
	puz2 := "Crossword-PreAndPost.puz"
	correctDecodeAndEncode(puz2, t)
}

func TestScrambleAndUnscramble(t *testing.T) {
	puz1 := "Crossword-Scrambled.puz"
	decodedBytes, err := loadFile(puz1)
	if err != nil {
		t.Errorf("Failed to load %s: %s", puz1, err)
	}

	puzzle, err := puz.DecodePuz(decodedBytes)
	if err != nil {
		t.Errorf("Failed to decode %s: %s", puz1, err)
	}

	if !puzzle.Scrambled() {
		t.Errorf("Puzzle %s should be marked as scrambled but is not", puz1)
	}

	err = puzzle.Unscramble(1111)
	if err == nil {
		t.Errorf("Puzzle %s should not have been unscrambled, wrong code used", puz1)
	}

	err = puzzle.Unscramble(1234)
	if err != nil {
		t.Errorf("Puzzle %s should have been unscrambled: %s", puz1, err)
	}

	puzCheck := "Crossword.puz"
	checkPuzzleBytes, err := loadFile(puzCheck)
	if err != nil {
		t.Errorf("Failed to load %s: %s", puzCheck, err)
	}

	checkPuzzle, err := puz.DecodePuz(checkPuzzleBytes)
	if err != nil {
		t.Errorf("Failed to decode %s: %s", puzCheck, err)
	}

	for y := range puzzle.Solution {
		if !bytes.Equal(puzzle.Solution[y], checkPuzzle.Solution[y]) {
			t.Errorf("Unscrambled solution does not match expected solution")
		}
	}

	//TODO: rescramble and make sure scrambled board matches loaded board
}

func correctDecodeAndEncode(name string, t *testing.T) {
	decodedBytes, err := loadFile(name)
	if err != nil {
		t.Errorf("Failed to load %s: %s", name, err)
	}

	puzzle, err := puz.DecodePuz(decodedBytes)
	if err != nil {
		t.Errorf("Failed to decode %s: %s", name, err)
	}

	encodedBytes, err := puz.EncodePuz(puzzle)
	if err != nil {
		t.Errorf("Failed to encode %s: %s", name, err)
	}

	if !bytes.Equal(decodedBytes, encodedBytes) {
		t.Errorf("Decoded and Encoded byte mismatch for %s", name)
	}
}

func loadFile(name string) ([]byte, error) {
	fp, err := os.Open(test_files + name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
