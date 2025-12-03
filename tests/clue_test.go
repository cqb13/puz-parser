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

func TestClueDirectionSorting(t *testing.T) {
	var clues puz.Clues

	clues = append(clues, puz.NewClue("clue 2", 1, 0, 0, puz.Down))
	clues = append(clues, puz.NewClue("clue 1", 1, 0, 0, puz.Across))

	clues.Sort()

	if clues[0].Direction != puz.Across {
		t.Fatalf("Failed to properly sort clues by direction at the same position")
	}
}

func TestCluePositionSorting(t *testing.T) {
	var clues puz.Clues

	clues = append(clues, puz.NewClue("clue 3", 1, 1, 0, puz.Across))
	clues = append(clues, puz.NewClue("clue 1", 0, 0, 0, puz.Across))
	clues = append(clues, puz.NewClue("clue 2", 1, 0, 0, puz.Across))

	clues.Sort()

	if clues[0].Clue != "clue 1" || clues[1].Clue != "clue 2" || clues[2].Clue != "clue 3" {
		t.Fatalf("Failed to properly sort clues by position")
	}
}

func TestAddingAndRemovingClues(t *testing.T) {
	p := puz.NewPuzzle(5, 5)

	if len(p.Clues()) != p.ExpectedClues() {
		t.Fatalf("Initial clue count did not match expected clue count")
	}

	if p.ExpectedClues() != 0 {
		t.Fatalf("No clues should be expected in a new puzzle")
	}

	p.Board[0][0].Value = puz.SOLID_SQUARE

	// adding clues
	ok := p.AddClue(puz.NewClue("clue 1", 1, 1, 0, puz.Down), false)
	if !ok {
		t.Fatalf("Failed to add valid clue 1,0.D")
	}

	if p.ExpectedClues() != 1 {
		t.Fatalf("Failed to increase expected clues on clue add")
	}

	ok = p.AddClue(puz.NewClue("clue 2", 1, 1, 0, puz.Across), false)
	if !ok {
		t.Fatalf("Failed to add valid clue 1,0.A")
	}

	if p.ExpectedClues() != 2 {
		t.Fatalf("Failed to increase expected clues on clue add")
	}

	ok = p.AddClue(puz.NewClue("clue -1", 1, 1, 0, puz.Across), false)
	if ok {
		t.Fatalf("Allowed placing a repeat clue 1,0.A")
	}

	if p.ExpectedClues() != 2 {
		t.Fatalf("Made change to expected clue count despite not allowing placement")
	}

	ok = p.AddClue(puz.NewClue("clue -1", 2, 1, 0, puz.Across), false)
	if ok {
		t.Fatalf("Allowed adding a clue at a position with a repeat number")
	}

	if p.ExpectedClues() != 2 {
		t.Fatalf("Made change to expected clue count despite not allowing placement")
	}

	ok = p.AddClue(puz.NewClue("clue 3", 2, 2, 2, puz.Across), false)
	if !ok {
		t.Fatalf("Prevented placing valid clue at invalid position when not validating")
	}

	if p.ExpectedClues() != 3 {
		t.Fatalf("Failed to increase expected clues on clue add")
	}

	ok = p.AddClue(puz.NewClue("clue -1", 2, 2, 2, puz.Down), true)
	if ok {
		t.Fatalf("Failed to prevent placing valid clue at invalid position when validating")
	}

	if p.ExpectedClues() != 3 {
		t.Fatalf("Made change to expected clue count despite not allowing placement")
	}

	ok = p.AddClue(puz.NewClue("clue 4", 3, 0, 1, puz.Down), true)
	if !ok {
		t.Fatalf("Prevented placing valid clue at valid position when validating")
	}

	if p.ExpectedClues() != 4 {
		t.Fatalf("Failed to increase expected clues on clue add")
	}

	ok = p.AddClue(puz.NewClue("clue -1", 4, 0, 4, puz.Down), true)
	if ok {
		t.Fatalf("Failed to prevent placing valid clue at valid position with invalid direction when validating")
	}

	if p.ExpectedClues() != 4 {
		t.Fatalf("Made change to expected clue count despite not allowing placement")
	}

	// removing clues
	p.RemoveClueByNum(4, puz.Down)
	if p.ExpectedClues() != 4 {
		t.Fatalf("Decreased expected clues when no clue was removed")
	}

	p.RemoveClueByNum(1, puz.Across)
	if p.ExpectedClues() != 3 {
		t.Fatalf("Failed to decrease expected clues when a clue was removed")
	}

	p.RemoveClueByPos(1, 0, puz.Down)
	if p.ExpectedClues() != 2 {
		t.Fatalf("Failed to decrease expected clues when a clue was removed")
	}

	if len(p.Clues()) != p.ExpectedClues() {
		t.Fatalf("Actual clues and expected clues became out of sync")
	}
}
