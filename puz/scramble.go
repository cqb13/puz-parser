package puz

import (
	"fmt"
	"unicode"
)

// a-z and A-Z
func isLetter(char byte) bool {
	if (char >= 0x41 && char <= 0x5A) || (char >= 0x61 && char <= 0x7A) {
		return true
	}

	return false
}

func createScrambleBuffer(puzzle *Puzzle) []byte {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()

	var buffer = make([]byte, width*height)
	var n = 0

	for x := range width {
		for y := range height {
			ch := puzzle.Board[y][x].Value
			if isLetter(ch) {
				buffer[n] = byte(unicode.ToUpper(rune(ch)))
				n++
			}
		}
	}

	return buffer[:n]
}

func updatePuzzleSolution(puzzle *Puzzle, buffer []byte) {
	height := puzzle.Board.Height()
	width := puzzle.Board.Width()
	var n = 0

	for x := range width {
		for y := range height {
			if isLetter(puzzle.Board[y][x].Value) {
				puzzle.Board[y][x].Value = buffer[n]
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

func convertLettersToNumbers(buffer []byte) {
	for i := range buffer {
		buffer[i] -= 'A'
	}
}

func convertNumbersToLetters(buffer []byte) {
	for i := range buffer {
		buffer[i] += 'A'
	}
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

func scramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	letterBuffer := createScrambleBuffer(puzzle)
	totalLetters := len(letterBuffer)

	if totalLetters < 12 {
		return ErrTooFewCharactersToScramble
	}

	puzzle.scramble.scrambledChecksum = checksumRegion(letterBuffer, 0)
	puzzle.scramble.scrambledTag = 4

	convertLettersToNumbers(letterBuffer)

	tempBuffer := make([]byte, totalLetters)

	copy(tempBuffer, letterBuffer)
	position := -1
	for i := range totalLetters {
		position += 16
		for position >= totalLetters {
			if totalLetters%2 == 0 {
				position -= totalLetters + 1
			} else {
				position -= totalLetters
			}
		}
		letterBuffer[position] = tempBuffer[i]
	}

	for round := range 4 {
		stepSize := 1 << (4 - round)

		position := -1
		for i := range totalLetters {
			position += stepSize
			for position >= totalLetters {
				if totalLetters%2 == 0 {
					position -= totalLetters + 1
				} else {
					position -= totalLetters
				}
			}
			keyOffset := int(keyDigits[i%4])
			letterBuffer[position] = byte((int(letterBuffer[position]) + keyOffset) % 26)
		}

		if stepSize > totalLetters {
			if totalLetters%2 == 0 {
				stepSize -= totalLetters + 1
			} else {
				stepSize -= totalLetters
			}
		}

		for repeat := 0; repeat < int(keyDigits[round]); repeat++ {
			copy(tempBuffer, letterBuffer[:stepSize])

			if totalLetters%2 == 0 {
				last := tempBuffer[stepSize-1]
				copy(tempBuffer[1:], tempBuffer[:stepSize-1])
				tempBuffer[0] = last
			}

			copy(letterBuffer, letterBuffer[stepSize:])
			copy(letterBuffer[totalLetters-stepSize:], tempBuffer[:stepSize])
		}
	}

	convertNumbersToLetters(letterBuffer)
	updatePuzzleSolution(puzzle, letterBuffer)

	return nil
}

func unscramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	letterBuffer := createScrambleBuffer(puzzle)
	totalLetters := len(letterBuffer)

	if totalLetters < 12 {
		return ErrTooFewCharactersToUnscramble
	}

	convertLettersToNumbers(letterBuffer)
	tempBuffer := make([]byte, totalLetters)

	for round := 3; round >= 0; round-- {
		stepSize := 1 << (4 - round)

		if stepSize > totalLetters {
			if stepSize%2 == 0 {
				stepSize -= totalLetters + 1
			} else {
				stepSize -= totalLetters
			}
		}

		for repeat := 0; repeat < int(keyDigits[round]); repeat++ {
			copy(tempBuffer, letterBuffer[totalLetters-stepSize:])

			if totalLetters%2 == 0 {
				first := tempBuffer[0]
				copy(tempBuffer, tempBuffer[1:stepSize])
				tempBuffer[stepSize-1] = first
			}

			copy(letterBuffer[stepSize:], letterBuffer[:totalLetters-stepSize])
			copy(letterBuffer[:stepSize], tempBuffer[:stepSize])
		}

		position := -1
		for i := range totalLetters {
			position += 1 << (4 - round)
			for position >= totalLetters {
				if totalLetters%2 == 0 {
					position -= totalLetters + 1
				} else {
					position -= totalLetters
				}
			}
			keyOffset := int(keyDigits[i%4])
			letterBuffer[position] = byte((int(letterBuffer[position]) - keyOffset + 26) % 26)
		}
	}

	copy(tempBuffer, letterBuffer)
	position := -1
	for i := range totalLetters {
		position += 16
		for position >= totalLetters {
			if totalLetters%2 == 0 {
				position -= totalLetters + 1
			} else {
				position -= totalLetters
			}
		}
		letterBuffer[i] = tempBuffer[position]
	}

	convertNumbersToLetters(letterBuffer)

	if checksumRegion(letterBuffer, 0) != puzzle.scramble.scrambledChecksum {
		return ErrIncorrectKeyProvided
	}

	updatePuzzleSolution(puzzle, letterBuffer)
	puzzle.scramble.scrambledTag = 0
	puzzle.scramble.scrambledChecksum = 0x0000

	return nil
}
