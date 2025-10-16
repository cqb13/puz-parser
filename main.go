package main

import (
	"dev/cqb13/puz-reader/puz"
	"io"
	"os"
)

func main() {
	fp, err := os.Open("test.puz")
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(fp)
	if err != nil {
		panic(err)
	}

	_, err = puz.LoadPuz(bytes)
	if err != nil {
		panic(err)
	}
}
