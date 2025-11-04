package puz

import "errors"

var (
	ErrOutOfBoundsRead              = errors.New("Out of bounds read")
	ErrOutOfBoundsWrite             = errors.New("Out of bounds write")
	ErrUnreadableData               = errors.New("Data does not appear to represent a puzzle")
	ErrMissingFileMagic             = errors.New("Failed to find ACROSS&DOWN in bytes")
	ErrGlobalChecksumMismatch       = errors.New("Global checksum mismatch")
	ErrCIBChecksumMismatch          = errors.New("CIB checksum mismatch")
	ErrMaskedLowChecksumMismatch    = errors.New("Masked Low checksum mismatch")
	ErrMaskedHighChecksumMismatch   = errors.New("Masked High checksum mismatch")
	ErrClueCountMismatch            = errors.New("The number of clues specified did not match the number of clues parsed")
	ErrExtraSectionChecksumMismatch = errors.New("Extra Section Checksum mismatch")
	ErrDuplicateExtraSection        = errors.New("A duplicate extra section was found")
	ErrUknownExtraSectionName       = errors.New("Unknown extra section name")
	ErrMissingExtraSection          = errors.New("An extra section was expected but not found")
	ErrBoardHeightMismatch          = errors.New("Board height did not match expected board height")
	ErrBoardWidthMismatch           = errors.New("Board width did not match expected board width")
)
