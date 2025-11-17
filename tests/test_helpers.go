package tests

import (
	"fmt"
	"io"
	"os"
	"strings"
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

func buildHex(b []byte) string {
	var out strings.Builder

	for i, v := range b {
		if i > 0 && i%16 == 0 {
			out.WriteString("\n")
		}

		if i%8 == 0 && i%16 != 0 {
			out.WriteString(" ")
		}

		fmt.Fprintf(&out, "%02x ", v)
	}

	return out.String()
}
