package puz

import (
	"fmt"
)

// a-z and A-Z
func isLetter(char byte) bool {
	if (char >= 0x41 && char <= 0x5A) || (char >= 0x61 && char <= 0x7A) {
		return true
	}

	return false
}

func createScramble(puzzle *Puzzle) string {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	scramble := ""

	for x := range width {
		for y := range height {
			ch := puzzle.Board[y][x].Value
			if isLetter(ch) {
				scramble += string(ch)
			}
		}
	}

	return scramble
}

func updatePuzzleSolution(puzzle *Puzzle, newSol string) {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	var n = 0

	for x := range width {
		for y := range height {
			if isLetter(puzzle.Board[y][x].Value) {
				puzzle.Board[y][x].Value = newSol[n]
				n++
			}
		}
	}
}

func keyToBytes(key int) ([]byte, error) {
	if key < 1000 || key > 9999 {
		return nil, ErrInvalidKeyLength
	}

	keyBytes := fmt.Appendf(nil, "%04d", key)

	for i := range keyBytes {
		if keyBytes[i] == '0' {
			return nil, ErrInvalidDigitInKey
		}
		keyBytes[i] -= '0'
	}

	return keyBytes, nil
}

func scramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	scramble := createScramble(puzzle)

	if len(scramble) < 12 {
		return ErrTooFewCharactersToScramble
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

func shiftString(unscrambled string, num int) string {
	return unscrambled[num:] + unscrambled[:num]
}

func scrambleString(unscrambled string) string {
	mid := len(unscrambled) / 2
	front := unscrambled[:mid]
	back := unscrambled[mid:]

	scrambled := ""

	for i := range len(front) {
		scrambled += string(back[i]) + string(front[i])
	}

	if len(unscrambled)%2 != 0 {
		scrambled += string(back[len(back)-1])
	}

	return scrambled
}

func unscramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	solution := createScramble(puzzle)
	if len(solution) < 12 {
		return ErrTooFewCharactersToUnscramble
	}

	for round := 3; round >= 0; round-- {
		digit := int(keyDigits[round])

		solution = unscrambleString(solution)
		solution = unshiftString(solution, digit)

		undo := ""
		for i, ch := range solution {
			letter := byte(ch) - keyDigits[i%4]
			if letter < 'A' {
				letter += 26
			}
			undo += string(letter)
		}
		solution = undo
	}

	if checksumRegion([]byte(solution), 0) != puzzle.scramble.scrambledChecksum {
		return ErrIncorrectKeyProvided
	}

	updatePuzzleSolution(puzzle, solution)
	puzzle.scramble.scrambledTag = 0
	puzzle.scramble.scrambledChecksum = 0x0000

	return nil
}

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

func unshiftString(s string, num int) string {
	num = num % len(s)
	return s[len(s)-num:] + s[:len(s)-num]
}

