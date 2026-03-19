package puz

import (
	"fmt"
	"strings"
)

// isLetter reports if char is either a-z or A-Z.
func isLetter(char byte) bool {
	if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') {
		return true
	}

	return false
}

// createScramble converts the crossword board into a string used in the scrambling/unscrambling algorithm.
//
// The string is formed by starting in the top left corner of the board and traversing down each column adding each letter and skipping any solid squares.
// returns NonLetterCharactersInScrambleError if the board contains a non letter character.
func createScramble(puzzle *Puzzle) (string, error) {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	var scramble strings.Builder

	for x := range width {
		for y := range height {
			ch := puzzle.Board[y][x].Answer
			if isLetter(ch) {
				scramble.WriteString(string(ch))
				continue
			}

			if ch != SolidSquare && ch != DiagramlessSolidSquare {
				return "", NonLetterCharactersInScrambleError
			}
		}
	}

	return scramble.String(), nil
}

// updatePuzzleSolution takes in the string used in scrambling process and puts the letters back in the board in the appropriate places.
func updatePuzzleSolution(puzzle *Puzzle, newSol string) {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	var n = 0

	for x := range width {
		for y := range height {
			if isLetter(puzzle.Board[y][x].Answer) {
				puzzle.Board[y][x].Answer = newSol[n]
				n++
			}
		}
	}
}

// keyToBytes returns an array of 4 bytes where each byte is a digit in the key.
//
// returns InvalidKeyLengthError if the key is less than 1000 or more than 9999.
// returns InvalidDigitInKeyError if the key contains a 0.
func keyToBytes(key int) ([]byte, error) {
	if key < 1000 || key > 9999 {
		return nil, InvalidKeyLengthError
	}

	keyBytes := fmt.Appendf(nil, "%04d", key)

	for i := range keyBytes {
		if keyBytes[i] == '0' {
			return nil, InvalidDigitInKeyError
		}
		keyBytes[i] -= '0'
	}

	return keyBytes, nil
}

// scramble scrambles a crosswords answer board.
//
// returns TooFewCharactersToScrambleError if the board has less than 12 characters. Also fails if the key is invalid or non letter characters are found in the board.
func scramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	scramble, err := createScramble(puzzle)
	if err != nil {
		return err
	}

	if len(scramble) < 12 {
		return TooFewCharactersToScrambleError
	}

	for _, digit := range keyDigits {
		lastScramble := scramble
		scramble = ""

		for i, letter := range lastScramble {
			letterVal := byte(letter) + keyDigits[i%4]

			// make sure letters are uppercase
			if letterVal > 90 {
				letterVal -= 26
			}

			scramble += string(letterVal)
		}

		scramble = shiftString(scramble, int(digit))
		scramble = scrambleString(scramble)
	}

	updatePuzzleSolution(puzzle, scramble)
	puzzle.scramble.scrambledTag = 4
	puzzle.scramble.scrambledChecksum = checksumRegion([]byte(scramble), 0)

	return nil
}

// TODO: doc this
func shiftString(unscrambled string, num int) string {
	return unscrambled[num:] + unscrambled[:num]
}

// TODO: doc this
func scrambleString(unscrambled string) string {
	mid := len(unscrambled) / 2
	front := unscrambled[:mid]
	back := unscrambled[mid:]

	var scrambled strings.Builder

	for i := range len(front) {
		scrambled.WriteString(string(back[i]) + string(front[i]))
	}

	if len(unscrambled)%2 != 0 {
		scrambled.WriteString(string(back[len(back)-1]))
	}

	return scrambled.String()
}

// unscramble unscrambles a crosswords board.
//
// returns TooFewCharactersToUncrambleError if the board has less than 12 characters. Also fails if the key is invalid or non letter characters are found in the board.
func unscramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	solution, err := createScramble(puzzle)
	if err != nil {
		return err
	}

	if len(solution) < 12 {
		return TooFewCharactersToUnscrambleError
	}

	for round := 3; round >= 0; round-- {
		digit := int(keyDigits[round])

		solution = unscrambleString(solution)
		solution = unshiftString(solution, digit)

		var undo strings.Builder
		for i, ch := range solution {
			letter := byte(ch) - keyDigits[i%4]
			if letter < 'A' {
				letter += 26
			}
			undo.WriteString(string(letter))
		}
		solution = undo.String()
	}

	if checksumRegion([]byte(solution), 0) != puzzle.scramble.scrambledChecksum {
		return IncorrectKeyProvidedError
	}

	updatePuzzleSolution(puzzle, solution)
	puzzle.scramble.scrambledTag = 0
	puzzle.scramble.scrambledChecksum = 0x0000

	return nil
}

// TODO: doc this
func unscrambleString(scrambled string) string {
	scrambledLen := len(scrambled)
	mid := scrambledLen / 2

	front := make([]byte, mid)
	back := make([]byte, scrambledLen-mid)

	j := 0
	for i := range mid {
		back[i] = scrambled[j]
		front[i] = scrambled[j+1]
		j += 2
	}

	if scrambledLen%2 != 0 {
		back[len(back)-1] = scrambled[len(scrambled)-1]
	}

	return string(front) + string(back)
}

// TODO: doc this
func unshiftString(s string, num int) string {
	num = num % len(s)
	return s[len(s)-num:] + s[:len(s)-num]
}
