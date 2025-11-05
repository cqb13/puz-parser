package puz

type Board [][]byte

func (board Board) inBounds(x int, y int) bool {
	if x < 0 || y < 0 || y >= len(board) || x >= len(board[y]) {
		return false
	}

	return true
}

func (board Board) isStartOfWord(x int, y int, dir Direction) bool {
	if !board.inBounds(x, y) || board[y][x] == BLACK_SQUARE {
		return false
	}

	height := len(board)
	width := len(board[y])

	// for across if cell is on edge or has black square to the left and has at least 1 cell to the right without a blacks square
	if dir == ACROSS && (x == 0 || (x-1 >= 0 && board[y][x-1] == BLACK_SQUARE)) && (x+1 < width && board[y][x+1] != BLACK_SQUARE) {
		return true
	}

	// for down if a cell is on edge or has a black square above it and has at least 1 cell bellow without a black square
	if dir == DOWN && (y == 0 || (y-1 >= 0 && board[y-1][x] == BLACK_SQUARE)) && (y+1 < height && board[y+1][x] != BLACK_SQUARE) {
		return true
	}

	return false
}

// Returns ok if a word of the given length can be placed at a position, not ok if the length of the wword is greater than the width/height of the board, or the word will cross a black square while placing
func (board Board) CanPlace(wordLen int, x int, y int, dir Direction) bool {
	if !board.inBounds(x, y) {
		return false
	}

	// make sure the word fits in bounds
	// -1 as letter will also be placed at origin
	if (dir == ACROSS && !board.inBounds(x+wordLen-1, y)) || (dir == DOWN && !board.inBounds(x, y+wordLen-1)) {
		return false
	}

	xOffset := x
	yOffset := y

	// make sure the word wont run into any black squares
	for range wordLen {
		if board[yOffset][xOffset] == BLACK_SQUARE {
			return false
		}

		if dir == ACROSS {
			xOffset++
		} else {
			yOffset++
		}
	}

	return true
}

// Places the word in the grid. Not ok if the length of the wword is greater than the width/height of the board, or the word will cross a black square while placing
func (board Board) PlaceWord(word string, x int, y int, dir Direction) bool {
	if !board.CanPlace(len(word), x, y, dir) {
		return false
	}

	xOffset := x
	yOffset := y

	for _, char := range []byte(word) {
		board[yOffset][xOffset] = char

		if dir == ACROSS {
			xOffset++
		} else {
			yOffset++
		}
	}

	return true
}

// Returns ok and the word if a word is found to start at the location x, y
func (board Board) GetWord(x int, y int, dir Direction) (string, bool) {
	if !board.isStartOfWord(x, y, dir) {
		return "", false
	}

	word := ""

	for {
		if y >= len(board) || x >= len(board[y]) || x < 0 || y < 0 {
			break
		}

		if board[y][x] == BLACK_SQUARE {
			break
		}

		word += string(board[y][x])

		if dir == ACROSS {
			x++
		} else {
			y++
		}
	}

	// a word must be at least 2 chars long
	if len(word) < min_word_len {
		return "", false
	}

	return word, true
}

func (board Board) GetWords(dir Direction) []string {
	var words []string
	for y := range len(board) {
		for x := range len(board[y]) {
			word, ok := board.GetWord(x, y, dir)
			if !ok {
				continue
			}

			words = append(words, word)
		}
	}

	return words
}
