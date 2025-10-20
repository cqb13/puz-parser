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
	var buffer = make([]byte, puzzle.Size)
	var n = 0

	for x := range puzzle.Width {
		for y := range puzzle.Height {
			ch := puzzle.Solution[y][x]
			if isLetter(ch) {
				buffer[n] = byte(unicode.ToUpper(rune(ch)))
				n++
			}
		}
	}

	return buffer[:n]
}

func updatePuzzleSolution(puzzle *Puzzle, buffer []byte) {
	var n = 0

	for x := range puzzle.Width {
		for y := range puzzle.Height {
			if isLetter(puzzle.Solution[y][x]) {
				puzzle.Solution[y][x] = buffer[n]
				n++
			}
		}
	}
}

func keyToBytes(key int) ([]byte, error) {
	if key < 1000 || key > 9999 {
		return nil, fmt.Errorf("the key must be a 4-digit number between 1000 and 9999")
	}

	keyBytes := fmt.Appendf(nil, "%04d", key)

	for i := range keyBytes {
		if keyBytes[i] == '0' {
			return nil, fmt.Errorf("the key cannot contain any zeros")
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

func Unscramble(puzzle *Puzzle, key int) error {
	keyDigits, err := keyToBytes(key)
	if err != nil {
		return err
	}

	letterBuffer := createScrambleBuffer(puzzle)
	totalLetters := len(letterBuffer)

	if totalLetters < 12 {
		return fmt.Errorf("too few characters to unscramble (minimum 12, found %d)", totalLetters)
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

	if checksumRegion(letterBuffer, 0) != puzzle.metadata.ScrambledChecksum {
		return fmt.Errorf("incorrect key provided (checksum mismatch)")
	}

	updatePuzzleSolution(puzzle, letterBuffer)
	puzzle.metadata.ScrambledTag = 0
	puzzle.metadata.ScrambledChecksum = 0x0000

	return nil
}
