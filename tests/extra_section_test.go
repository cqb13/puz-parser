package tests

import (
	"slices"
	"testing"

	"github.com/cqb13/puz-parser/puz"
)

func TestGRBSandRTBL(t *testing.T) {
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

func TestAllSections(t *testing.T) {
	puzzles := []string{
		"All-Sections-Sorted.puz",
		"All-Sections-Unsorted.puz",
	}

	for _, name := range puzzles {
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

		if !puzzle.HasExtraSection(puz.TimerSection) {
			t.Errorf("Failed to find expected LTIM section")
		}

		if !puzzle.HasExtraSection(puz.MarkupBoardSection) {
			t.Errorf("Failed to find expected GEXT section")
		}

		if !puzzle.HasExtraSection(puz.UserRebusTableSection) {
			t.Errorf("Failed to find expected RUSR section")
		}
	}
}

func TestAddingAndRemoving(t *testing.T) {
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

	ok := puzzle.RemoveExtraSection(puz.RebusSection)
	if !ok {
		t.Errorf("Failed to remove GRBS section")
	}

	if puzzle.HasExtraSection(puz.RebusSection) {
		t.Errorf("GRBS section present after removal")
	}

	ok = puzzle.AddExtraSection(puz.RebusSection)

	if !ok {
		t.Errorf("Failed to add GRBS section")
	}

	if !puzzle.HasExtraSection(puz.RebusSection) {
		t.Errorf("Failed to find added GRBS section")
	}
}

func TestExtraSectionSorting(t *testing.T) {
	sortedName := "All-Sections-Sorted.puz"
	unsortedName := "All-Sections-Unsorted.puz"

	sortedData, err := loadFile(sortedName)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", sortedName, err)
	}

	unsortedData, err := loadFile(unsortedName)
	if err != nil {
		t.Fatalf("Failed to load %s: %v", unsortedName, err)
	}

	if slices.Compare(sortedData, unsortedData) == 0 {
		t.Fatalf("Sorted and unsorted section puzzles were the same")
	}

	puzzle, err := puz.DecodePuz(unsortedData)
	if err != nil {
		t.Fatalf("Failed to decode %s: %v", unsortedName, err)
	}

	puzzle.SortExtraSections()

	encoded, err := puz.EncodePuz(puzzle)
	if err != nil {
		t.Fatalf("Failed to encode %s: %v", unsortedName, err)
	}

	if slices.Compare(sortedData, encoded) != 0 {
		t.Fatalf("Sorted section order did not match expected sort")
	}
}
