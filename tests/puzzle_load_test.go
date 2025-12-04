package tests

import (
	"bytes"
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

// ensure data (name, author, notes, version, etc...)  is loaded correctly
func TestPuzzleLoading(t *testing.T) {
	name := "Crossword.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	p, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	if p.Title != "2025!" {
		t.Fatalf("Found unexpected title, expected 2025!, found %s", p.Title)
	}

	if p.Author != "cqb13" {
		t.Fatalf("Found unexpected author, expected cqb13, found %s", p.Author)
	}

	if !bytes.Equal([]byte(p.Copyright), []byte{0x54, 0x61, 0x6C, 0x6F, 0x6E, 0x20, 0x47, 0x61, 0x6D, 0x65, 0x73, 0x20, 0xA9, 0x20, 0x32, 0x30, 0x32, 0x35}) {
		t.Fatalf("Found unexpected copyright, expected 'Talon Games © 2025', found '%s'", p.Copyright)
	}

	if p.Notes != "https://games.shstalon.com/games/crossword?type=mini" {
		t.Fatalf("Found unexpected note, expected 'https://games.shstalon.com/games/crossword?type=mini', found %s", p.Notes)
	}

	if p.Version() != "1.4" {
		t.Fatalf("Found unexpected version, expected 1.4, found %s", p.Version())
	}
}
