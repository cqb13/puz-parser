package puz_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func loadFile(t *testing.T, name string) []byte {
	t.Helper()

	data, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("failed to load %s: %v", name, err)
	}

	return data
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
