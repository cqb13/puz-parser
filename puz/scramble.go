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

func createSolutionBuffer(puzzle *Puzzle) []byte {
	var buffer = make([]byte, puzzle.Width*puzzle.Height)
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

func Unscramble(puzzle *Puzzle, key int) error {
	keyBytes, err := keyToBytes(key)
	if err != nil {
		return err
	}

	buffer := createSolutionBuffer(puzzle)
	size := len(buffer)

	if size < 12 {
		return fmt.Errorf("Too few characters to unscramble, minimum 12, found %d", size)
	}

	for i, b := range buffer {
		buffer[i] = b - 'A'
	}

	tmp := make([]byte, size)

	for k := 3; k >= 0; k-- {
		n := 1 << (4 - k)
		if n > size {
			n -= size | 1
		}
		for i := 0; i < int(keyBytes[k]); i++ {
			copy(tmp, buffer[size-n:])

			if size%2 == 0 {
				tmp = append(tmp[1:n], tmp[0])
			}

			copy(buffer[n:], buffer[:size-n])
			copy(buffer[:n], tmp[:n])
		}

		j := -1
		for i := range size {
			j += 1 << (4 - k)
			for j >= size {
				j -= size | 1
			}
			buffer[j] = byte((int(buffer[j]) - int(keyBytes[i%4]) + 26) % 26)
		}
	}

	copy(tmp, buffer)
	j := -1
	for i := range size {
		j += 16
		for j >= size {
			j -= size | 1
		}
		buffer[i] = tmp[j]
	}

	for i := range buffer {
		buffer[i] += 'A'
	}

	if checksumRegion(buffer, 0) != puzzle.metadata.ScrambledChecksum {
		return fmt.Errorf("Incorrect key provided, checksum mismatch")
	}

	updatePuzzleSolution(puzzle, buffer)
	puzzle.metadata.ScrambledTag = 0
	puzzle.metadata.ScrambledChecksum = 0x0000
	return nil
}
