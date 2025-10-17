package main

import (
	"dev/cqb13/puz-reader/puz"
	"fmt"
	"io"
	"os"
)

func main() {
	var basePath = "./tests/test-files/"
	entries, err := os.ReadDir(basePath)
	if err != nil {
		panic(err)
	}

	analyzed := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fp, err := os.Open(basePath + entry.Name())
		if err != nil {
			fmt.Printf("Failed to open %s: %s\n", entry.Name(), err)
			continue
		}
		defer fp.Close()

		bytes, err := io.ReadAll(fp)
		if err != nil {
			fmt.Printf("Failed to read %s: %s\n", entry.Name(), err)
			continue
		}

		_, err = puz.DecodePuz(bytes)
		if err != nil {
			fmt.Printf("Failed to parse %s: %s\n", entry.Name(), err)
			continue
		}
		analyzed++
	}

	fmt.Printf("Analyzed %d/%d\n", analyzed, len(entries))
}
