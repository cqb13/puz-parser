package puz

type Board [][]Cell

type Word struct {
	Word      string
	Num       int
	StartX    int
	StartY    int
	Direction Direction
}

type Cell struct {
	Value    byte
	State    byte
	RebusKey byte
	Markup   byte
}

func NewBoard(width uint8, height uint8) [][]Cell {
	board := make([][]Cell, height)

	for y := range height {
		board[y] = make([]Cell, width)

		for x := range width {
			board[y][x] = Cell{
				EMPTY_SOLUTION_SQUARE,
				EMPTY_STATE_SQUARE,
				0x00,
				0x00,
			}
		}
	}

	return board
}

func NewBoardFromArr(byteBoard [][]byte) ([][]Cell, error) {
	board := make([][]Cell, len(byteBoard))

	prevWdith := len(byteBoard[0])
	for y, row := range byteBoard {
		board[y] = make([]Cell, len(row))
		if len(row) != prevWdith {
			return nil, ErrBoardWidthMismatch
		}
		for x, value := range row {
			cell := Cell{
				EMPTY_SOLUTION_SQUARE,
				EMPTY_STATE_SQUARE,
				0x00,
				0x00,
			}

			if value == SOLID_SQUARE || value == DIAGRAMLESS_SOLID_SQUARE {
				cell.State = value
			}

			cell.Value = value

			board[y][x] = cell
		}
	}

	return board, nil
}

func (b Board) Height() int {
	return len(b)
}

func (b Board) Width() int {
	return len(b[0])
}

func (b Board) inBounds(x int, y int) bool {
	if x >= 0 && x < b.Width() && y >= 0 && y < b.Height() {
		return true
	}

	return false
}

func (b Board) IsSolidSquare(x int, y int) bool {
	return b[y][x].Value == SOLID_SQUARE || b[y][x].Value == DIAGRAMLESS_SOLID_SQUARE
}

func (b Board) GetWord(x int, y int, dir Direction) (string, bool) {
	if !b.inBounds(x, y) {
		return "", false
	}

	if b.IsSolidSquare(x, y) {
		return "", false
	}

	word := ""

	xOffset := x
	yOffset := y

	for {
		word += string(b[yOffset][xOffset].Value)

		if dir == Across {
			xOffset += 1
		} else {
			yOffset += 1
		}

		if !b.inBounds(xOffset, yOffset) || b.IsSolidSquare(xOffset, yOffset) {
			break
		}
	}

	return word, true
}

func (b Board) StartsAcrossWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if b.IsSolidSquare(x, y) {
		return false
	}

	if x == 0 || b.IsSolidSquare(x-1, y) {
		if x+1 < b.Width() && !b.IsSolidSquare(x+1, y) {
			return true
		}
	}

	return false
}

func (b Board) StartsDownWord(x int, y int) bool {
	if !b.inBounds(x, y) {
		return false
	}

	if b.IsSolidSquare(x, y) {
		return false
	}

	if y == 0 || b.IsSolidSquare(x, y-1) {
		if y+1 < b.Height() && !b.IsSolidSquare(x, y+1) {
			return true
		}
	}

	return false
}

func (b Board) GetWords() []Word {
	var words []Word

	width := b.Width()
	nextWordNum := 1
	for y := range b.Height() {
		for x := range width {
			if b.IsSolidSquare(x, y) {
				continue
			}

			startsAcrossWord := b.StartsAcrossWord(x, y)
			startsDownWord := b.StartsDownWord(x, y)

			if startsAcrossWord {
				word, ok := b.GetWord(x, y, Across)
				if ok {
					words = append(words, Word{
						word,
						nextWordNum,
						x,
						y,
						Across,
					})
				}
			}

			if startsDownWord {
				word, ok := b.GetWord(x, y, Down)
				if ok {
					words = append(words, Word{
						word,
						nextWordNum,
						x,
						y,
						Down,
					})
				}
			}

			if startsAcrossWord || startsDownWord {
				nextWordNum++
			}
		}
	}

	return words
}
