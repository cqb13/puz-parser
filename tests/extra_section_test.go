package tests

import (
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestGRBSandRTBL(t *testing.T) {
	// puzzle should have a GRBS and RTBL with everything else being nil
	name := "Crossword-EXT-Rebus.puz"
	data, err := loadFile(name)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", name, err)
	}

	puzzle, err := puz.DecodePuz(data)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", name, err)
	}

	//TODO: add checks for other 2

	if puzzle.Extras.RebusTable == nil {
		t.Errorf("Failed to find expected RTBL section")
	}

	if puzzle.Extras.Timer != nil {
		t.Errorf("Found unexpected LTIM section")
	}

	if puzzle.Extras.UserRebusTable != nil {
		t.Errorf("Found unexpected RUSR section")
	}
}

//TODO: test puzzle with all sections
