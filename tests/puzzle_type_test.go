package tests

import (
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestDiagramlessPuzzle(t *testing.T) {
	name := "NYT-Diagramless.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	if puzzle.PuzzleType != puz.Diagramless {
		t.Fatalf("Expected %s to be a diagramless puzzle", name)
	}
}

func TestNormalPuzzle(t *testing.T) {
	name := "Crossword.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	if puzzle.PuzzleType != puz.Normal {
		t.Fatalf("Expected %s to be a diagramless puzzle", name)
	}
}
