package parser

import (
	"fmt"
	"os"
	"path/filepath"
)

// NewForFile returns the appropriate parser for a given lockfile name
func NewForFile(filename string) (LockfileParser, error) {
	base := filepath.Base(filename)
	switch base {
	case "package-lock.json":
		return NewNpmParser(), nil
	case "requirements.txt":
		return NewPipParser(), nil
	default:
		return nil, fmt.Errorf("unsupported lockfile: %s", base)
	}
}

// Detect looks for supported lockfiles in the given directory
func Detect(dir string) ([]string, error) {
	supported := []string{"package-lock.json", "requirements.txt"}
	found := []string{}

	for _, name := range supported {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			found = append(found, path)
		}
	}

	return found, nil
}
