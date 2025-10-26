package tests

import (
	"io"
	"os"
)

const testFilesDir = "./test-files/"

func loadFile(name string) ([]byte, error) {
	fp, err := os.Open(testFilesDir + name)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return io.ReadAll(fp)
}
