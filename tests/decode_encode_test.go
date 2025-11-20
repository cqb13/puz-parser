package tests

import (
	"bytes"
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestDecodeAndEncode(t *testing.T) {
	testCases := []string{
		"Crossword-Blank-Squares.puz",
		"Crossword-EXT-Rebus.puz",
		"Crossword-PreAndPost-Scrambled.puz",
		"Crossword-PreAndPost.puz",
		"Crossword-PreAndPost.puz",
		"Crossword-1.3.puz",
		"Crossword.puz",
		"washpost.puz",
	}

	for _, name := range testCases {
		t.Run(name, func(t *testing.T) {
			correctDecodeAndEncode(name, t)
		})
	}
}

func correctDecodeAndEncode(name string, t *testing.T) {
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	encoded, err := puz.EncodePuz(puzzle)
	if err != nil {
		t.Fatalf("Failed to encode %s: %v", name, err)
	}

	if !bytes.Equal(data, encoded) {
		t.Errorf("Encoded bytes do not match original for %s\n\noriginal:\n%s\n\nnew:\n%s", name, buildHex(data), buildHex(encoded))
	}
}
