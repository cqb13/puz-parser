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

	if !puzzle.HasExtraSection(puz.RebusSection) {
		t.Errorf("Failed to find expected GRBS section")
	}

	if !puzzle.HasExtraSection(puz.RebusTableSection) {
		t.Errorf("Failed to find expected RTBL section")
	}

	if puzzle.HasExtraSection(puz.TimerSection) {
		t.Errorf("Found unexpected LTIM section")
	}

	if puzzle.HasExtraSection(puz.MarkupBoardSection) {
		t.Errorf("Found unexpected GEXT section")
	}

	if puzzle.HasExtraSection(puz.UserRebusTableSection) {
		t.Errorf("Found unexpected RUSR section")
	}
}

//TODO: test puzzle with all sections

//TODO: test adding and removing sections

//TODO: test sorting sections
