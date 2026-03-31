package parser

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type pyproject struct {
	Project struct {
		Name    string      `toml:"name"`
		License interface{} `toml:"license"`
	} `toml:"project"`
}

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

	projectName := "python-project"
	projectLicense := "UNKNOWN"

	// Try to find pyproject.toml in the same directory
	dir := filepath.Dir(filePath)
	pyprojectPath := filepath.Join(dir, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectPath); err == nil {
		var config pyproject
		if err := toml.Unmarshal(data, &config); err == nil {
			if config.Project.Name != "" {
				projectName = config.Project.Name
			}

			// Handle license as string or table (PEP 621)
			switch l := config.Project.License.(type) {
			case string:
				projectLicense = NormalizeLicense(l)
			case map[string]interface{}:
				if text, ok := l["text"].(string); ok {
					projectLicense = NormalizeLicense(text)
				} else if file, ok := l["file"].(string); ok {
					// Fallback to filename if text is not provided
					projectLicense = "FILE:" + file
				}
			}
		}
	}

	rootBuilder := NewDependencyBuilder().
		Name(projectName).
		License(projectLicense).
		Ecosystem("pip").
		Depth(0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 1. Remove comments
		if idx := strings.Index(line, "#"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		// 2. Remove environment markers (anything after ;)
		spec := line
		if idx := strings.Index(line, ";"); idx != -1 {
			spec = strings.TrimSpace(line[:idx])
		}

		// 3. Extract name and version
		name := spec
		version := "latest"

		// Handle common specifiers: ==, >=, <=, !=, ~=, >, <
		operators := []string{"==", ">=", "<=", "!=", "~=", ">", "<"}
		var foundOp string
		var opIdx int = -1

		for _, op := range operators {
			if idx := strings.Index(spec, op); idx != -1 {
				if opIdx == -1 || idx < opIdx {
					opIdx = idx
					foundOp = op
				}
			}
		}

		if opIdx != -1 {
			name = strings.TrimSpace(spec[:opIdx])
			version = strings.TrimSpace(spec[opIdx+len(foundOp):])
			// If there are multiple specifiers (e.g. pkg>=1.0,<2.0), just take the first part
			if idx := strings.Index(version, ","); idx != -1 {
				version = strings.TrimSpace(version[:idx])
			}
		}

		dep, err := NewDependencyBuilder().
			Name(name).
			Version(version).
			License("UNKNOWN"). // Will be enriched later via PyPI API
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
