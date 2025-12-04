package tests

import (
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestSolidSquareDetection(t *testing.T) {
	var board puz.Board = puz.NewBoard(5, 5)

	board[0][0].Value = puz.SOLID_SQUARE
	board[0][1].Value = puz.DIAGRAMLESS_SOLID_SQUARE

	if !board.IsSolidSquare(0, 0) {
		t.Fatalf("Failed to detect solid square")
	}

	if !board.IsSolidSquare(1, 0) {
		t.Fatalf("Failed to detect diagramless solid square")
	}
}

func TestWordStartDetection(t *testing.T) {
	var board puz.Board = puz.NewBoard(5, 5)

	board[0][0].Value = puz.SOLID_SQUARE

	// make sure words cant start in solid squares
	if board.StartsAcrossWord(0, 0) {
		t.Fatalf("Detected the start of an across word in a solid square")
	}

	if board.StartsDownWord(0, 0) {
		t.Fatalf("Detected the start of a down word in a solid square")
	}

	// make sure words cant start without and edge or solid square
	if board.StartsAcrossWord(2, 2) {
		t.Fatalf("Detected the start of an across word without and edge or solid square")
	}

	if board.StartsDownWord(2, 2) {
		t.Fatalf("Detected the start of a down word without and edge or solid square")
	}

	// make sure words are detected on edges
	if !board.StartsAcrossWord(0, 2) {
		t.Fatalf("Failed to detect the start of an across word on an edge")
	}

	if !board.StartsDownWord(2, 0) {
		t.Fatalf("Failed to detect the start of a down word on an edge")
	}

	// make sure words are detected on solid square
	if !board.StartsAcrossWord(0, 1) {
		t.Fatalf("Failed to detect the start of an across word on an solid square")
	}

	if !board.StartsDownWord(1, 0) {
		t.Fatalf("Failed to detect the start of a down word on an solid square")
	}

	// make sure words aren't detected going off the board
	if board.StartsAcrossWord(4, 0) {
		t.Fatalf("Detected the start of an across word without enough space on the board")
	}

	if board.StartsDownWord(0, 4) {
		t.Fatalf("Detected the start of a down word without enough space on the board")
	}

	// make sure 2 letters long works
	board[2][2].Value = puz.SOLID_SQUARE

	if !board.StartsAcrossWord(0, 2) {
		t.Fatalf("Failed to detect the start of a 2 letter long across word")
	}

	if !board.StartsDownWord(2, 3) {
		t.Fatalf("Failed to detect the start of a 2 letter long down word")
	}

	// make sure 1 letter long words fail
	board[2][1].Value = puz.SOLID_SQUARE
	board[4][2].Value = puz.SOLID_SQUARE

	if board.StartsAcrossWord(0, 2) {
		t.Fatalf("Detected the start of a 1 letter long across word")
	}

	if board.StartsDownWord(2, 3) {
		t.Fatalf("Detected the start of a 1 letter long down word")
	}
}

func TestGettingWords(t *testing.T) {
	name := "Crossword.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	board := puzzle.Board

	// ensure the right amount of words are found with the correct directions
	words := board.GetWords()
	if len(words) != 10 {
		t.Fatalf("Failed to find expected amount of words in %s, expected 10, found %d", name, len(words))
	}

	acrossWords := 0
	for _, word := range words {
		if word.Direction == puz.Across {
			acrossWords++
		}
	}

	if acrossWords != 5 {
		t.Fatalf("Failed to properly determine the directions of all words in %s, expected 5 across words, found %d", name, acrossWords)
	}

	// ensure the words can be retrieved from the board
	word, ok := board.GetWord(0, 0, puz.Across)
	if !ok {
		t.Fatalf("Failed to retrieve word from valid position 0,0.A")
	}

	if word != "BASS" {
		t.Fatalf("Failed to retrieve correct word from 0,0.A, expected BASS, found %s", word)
	}

	word, ok = board.GetWord(0, 0, puz.Down)
	if !ok {
		t.Fatalf("Failed to retrieve word from valid position 0,0.D")
	}

	if word != "BASH" {
		t.Fatalf("Failed to retrieve correct word from 0,0.D, expected BASH, found %s", word)
	}

	word, ok = board.GetWord(1, 4, puz.Across)
	if !ok {
		t.Fatalf("Failed to retrieve word from valid position 1,4.A")
	}

	if word != "REED" {
		t.Fatalf("Failed to retrieve correct word from 1,4.A, expected REED, found %s", word)
	}

	word, ok = board.GetWord(4, 1, puz.Down)
	if !ok {
		t.Fatalf("Failed to retrieve word from valid position 4,1.D")
	}

	if word != "DEED" {
		t.Fatalf("Failed to retrieve correct word from 4,1.D, expected DEED, found %s", word)
	}
}
