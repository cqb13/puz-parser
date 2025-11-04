# Solution Board

needs a cell with corresponding clue # (0 for none) and letter

getWord(x, y, dir) str, ok (not ok if 0 clue #)

getWords(dir) []str all words in that direction

put(byte) sets a value in the grid

toBytes() [][]byte

# Clues

each clue should have the clue itself, the direction of the clue word and the word number it corresponds to

toBytes() [][]byte **remember to encode in proper order**

addClue(clue, direction, #) **might need to validate to ensure the clue can fit on board (check how sites respond)**

removeClue(clue, direction, #)
