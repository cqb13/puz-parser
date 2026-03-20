package puz_test

import (
	puz "github.com/cqb13/puz-parser"
	"testing"
)

type scrambleTestCase struct {
	name          string
	scrambledFile string
	plainFile     string
}

var testCases = []scrambleTestCase{
	{
		name:          "Crossword",
		scrambledFile: "Crossword-Scrambled.puz",
		plainFile:     "Crossword.puz",
	},
	{
		name:          "Crossword-PreAndPost",
		scrambledFile: "Crossword-PreAndPost-Scrambled.puz",
		plainFile:     "Crossword-PreAndPost.puz",
	},
}

func TestUnscramble(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodedBytes := loadFile(t, tc.scrambledFile)
			puzzle, err := puz.DecodePuz(decodedBytes)
			if err != nil {
				t.Fatalf("Failed to decode %s: %v", tc.scrambledFile, err)
			}

			if !puzzle.Scrambled() {
				t.Errorf("Puzzle %s should be marked as scrambled", tc.scrambledFile)
			}

			// Wrong code should fail
			if err := puzzle.Unscramble(1111); err == nil {
				t.Errorf("Puzzle %s should not unscramble with wrong code", tc.scrambledFile)
			}

			// Correct code should succeed
			if err := puzzle.Unscramble(1234); err != nil {
				t.Errorf("Puzzle %s failed to unscramble: %v", tc.scrambledFile, err)
			}

			// Compare with original unscrambled puzzle
			checkBytes := loadFile(t, tc.plainFile)

			checkPuzzle, err := puz.DecodePuz(checkBytes)
			if err != nil {
				t.Fatalf("Failed to decode %s: %v", tc.plainFile, err)
			}

			for y := range puzzle.Board {
				for x := range puzzle.Board[y] {
					if puzzle.Board[y][x].Answer != checkPuzzle.Board[y][x].Answer {
						t.Errorf("Cell x: %d y: %d mismatch after unscramble (%c != %c)", x, y, puzzle.Board[y][x].Answer, checkPuzzle.Board[y][x].Answer)

					}
				}
			}
		})
	}
}

func TestScramble(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decodedBytes := loadFile(t, tc.plainFile)

			puzzle, err := puz.DecodePuz(decodedBytes)
			if err != nil {
				t.Fatalf("Failed to decode %s: %v", tc.plainFile, err)
			}

			if puzzle.Scrambled() {
				t.Errorf("Puzzle %s should be marked as unscrambled", tc.plainFile)
			}

			// Scramble the puzzle
			if err := puzzle.Scramble(1234); err != nil {
				t.Errorf("Puzzle %s failed to scramble: %v", tc.plainFile, err)
			}

			// Compare with scrambled file
			checkBytes := loadFile(t, tc.scrambledFile)

			checkPuzzle, err := puz.DecodePuz(checkBytes)
			if err != nil {
				t.Fatalf("Failed to decode %s: %v", tc.scrambledFile, err)
			}

			for y := range puzzle.Board {
				for x := range puzzle.Board[y] {
					if puzzle.Board[y][x].Answer != checkPuzzle.Board[y][x].Answer {
						t.Errorf("Cell x: %d y: %d mismatch after scramble (%c != %c)", x, y, puzzle.Board[y][x].Answer, checkPuzzle.Board[y][x].Answer)
					}
				}
			}
		})
	}
}
