package tests

import (
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestClueLoading(t *testing.T) {
	name := "Crossword.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	clues := puzzle.Clues()

	if len(clues) != puzzle.ExpectedClues() {
		t.Fatalf("Clue amount did not match expected clue count on initial load")
	}
}

/*
Expected clues for this crossword

1.A: Lowest vocal range
1.D: Shell used for Unix commands
2.D: Person playing a role
3.D: Seaside land
4.D: Smell, touch, e.g.
5.A: Felt dull pain
6.D: A legal document that transfers ownership of property from one person to another
7.A: A pebble is a small one
8.A: A pony is a small one
9.A: A plant used for thatching
*/
func TestGettingClues(t *testing.T) {
	name := "Crossword.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	// getting clues by number
	clue, ok := puzzle.GetClueByNum(1, puz.Down)
	if !ok {
		t.Fatalf("Failed to get expected clue by 1.D")
	}

	if clue.Clue != "Shell used for Unix commands" {
		t.Fatalf("Clue 1.D did not match expected clue, expected: Shell used for Unix commands, found: %s", clue.Clue)
	}

	clue, ok = puzzle.GetClueByNum(7, puz.Across)
	if !ok {
		t.Fatalf("Failed to get expected clue by 7.A")
	}

	if clue.Clue != "A pebble is a small one" {
		t.Fatalf("Clue 7.A did not match expected clue, expected: A pebble is a small one, found: %s", clue.Clue)
	}

	_, ok = puzzle.GetClueByNum(111, puz.Across)
	if ok {
		t.Fatalf("Got a clue by a number that does not exist, 111.A")
	}

	// getting clues by pos
	clue, ok = puzzle.GetClueByPos(1, 4, puz.Across)
	if !ok {
		t.Fatalf("Failed to get expected clue by 1,4.A")
	}

	if clue.Clue != "A plant used for thatching" {
		t.Fatalf("Clue 1,4.A did not match expected clue, expected: A plant used for thatching, found: %s", clue.Clue)
	}

	clue, ok = puzzle.GetClueByPos(3, 0, puz.Down)
	if !ok {
		t.Fatalf("Failed to get expected clue by 0,3.D")
	}

	if clue.Clue != "Smell, touch, e.g." {
		t.Fatalf("Clue 0,3.D did not match expected clue, expected: Smell, touch, e.g., found: %s", clue.Clue)
	}

	_, ok = puzzle.GetClueByPos(1, 4, puz.Down)
	if ok {
		t.Fatalf("Got a clue by a pos that does not exist, 1,4.D")
	}
}

//TODO: tests adding and removing clues

//TODO: test sorting clues
