package parser

import (
	"bufio"
	"os"
	"strings"
)

type PipParser struct {
	BaseParser
}

func NewPipParser() *PipParser {
	p := &PipParser{}
	// Pip doesn't use the BaseParser.Decoder pattern easily because it's line-based
	// but we'll adapt it or just override Parse
	return p
}

// Parse pip requirements file (usually requirements.txt)
// See https://pip.pypa.io/en/stable/reference/requirements-file-format/#structure for details.
func (p *PipParser) Parse(filePath string) (*Dependency, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	rootBuilder := NewDependencyBuilder().
		Name("python-project").
		Ecosystem("pip").
		Depth(0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle name==version or name>=version
		name := line
		version := "latest"
		if strings.Contains(line, "==") {
			parts := strings.Split(line, "==")
			name = strings.TrimSpace(parts[0])
			version = strings.TrimSpace(parts[1])
		} else if strings.Contains(line, ">=") {
			parts := strings.Split(line, ">=")
			name = strings.TrimSpace(parts[0])
			version = strings.TrimSpace(parts[1])
		}

		dep, err := NewDependencyBuilder().
			Name(name).
			Version(version).
			License("UNKNOWN"). // Will be enriched later (maybe via PyPI API)
			Ecosystem("pip").
			Depth(1).
			Build()

		if err == nil {
			rootBuilder.AddDependency(dep)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return rootBuilder.Build()
}
