package puz

import (
	"errors"
	"fmt"
)

var (
	OutOfBoundsReadError               = errors.New("Out of bounds read")
	OutOfBoundsWriteError              = errors.New("Out of bounds write")
	UnreadableDataError                = errors.New("Data does not appear to represent a crossword puzzle")
	MissingFileMagicError              = errors.New("Failed to find ACROSS&DOWN string in file")
	UknownExtraSectionNameError        = errors.New("Unknown extra section name")
	MissingExtraSectionError           = errors.New("An extra section was expected but not found")
	BoardWidthMismatchError            = errors.New("Board contains rows of unequal width")
	PuzzleIsUnscrambledError           = errors.New("Puzzle is already unscrambled")
	PuzzleIsScrambledError             = errors.New("Puzzle is already scrambled")
	InvalidVersionFormatError          = errors.New("Invalid version format, must be X.X")
	TooFewCharactersToUnscrambleError  = errors.New("Too few characters to unscramble, minimum 12")
	TooFewCharactersToScrambleError    = errors.New("Too few characters to scramble, minimum 12")
	NonLetterCharactersInScrambleError = errors.New("Scramble operations can not be performed on grids with non-letter characters")
	InvalidDigitInKeyError             = errors.New("Key cannot contain any zeros")
	InvalidKeyLengthError              = errors.New("Key must be a 4-digit number")
	IncorrectKeyProvidedError          = errors.New("Failed to unscramble, incorrect key provided")
)

// Checksum Mismatch
type checksum int

const (
	globalChecksum checksum = iota
	cibChecksum
	maskedLowChecksum
	maskedHighChecksum
)

var checksumStrMap = map[checksum]string{
	globalChecksum:     "Global Checksum",
	cibChecksum:        "CIB Checksum",
	maskedLowChecksum:  "Masked Low Checksum",
	maskedHighChecksum: "Masked High Checksum",
}

func (c checksum) String() string {
	return checksumStrMap[c]
}

type ChecksumMismatchError struct {
	expected   int
	calculated int
	checksum   checksum
}

func (e *ChecksumMismatchError) Error() string {
	return fmt.Sprintf("%s mismatch: expected %d, calculated %d", e.checksum.String(), e.expected, e.calculated)
}

// Extra Section Checksum Mismatch
type ExtraSectionChecksumMismatchError struct {
	expected   uint16
	calculated uint16
	section    ExtraSection
}

func (e *ExtraSectionChecksumMismatchError) Error() string {
	return fmt.Sprintf("%s section checksum mismatch: expected %d, calculated %d", e.section.String(), e.expected, e.calculated)
}

// Clue Mismatch
type ClueCountMismatchError struct {
	expected int
	found    int
}

func (e *ClueCountMismatchError) Error() string {
	return fmt.Sprintf("The expected clue count did not match the number of clue: expected %d, found %d", e.expected, e.found)
}

// Duplicate Extra Section
type DuplicateExtraSectionError struct {
	section ExtraSection
}

func (e *DuplicateExtraSectionError) Error() string {
	return fmt.Sprintf("A duplicate %s section was found", e.section.String())
}
