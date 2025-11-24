package main

import (
	"fmt"
	"io"
	"os"

	"github.com/cqb13/puz-parser/puz"
)

func main() {
	fp, _ := os.Open("./tests/test-files/NYT-Diagramless.puz")
	defer fp.Close()

	bytes, _ := io.ReadAll(fp)

	p, _ := puz.DecodePuz(bytes)

	if p.PuzzleType == puz.Diagramless {
		fmt.Println("Diagramless")
	} else {
		fmt.Println("Normal")
	}
}
