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
		Name    string `toml:"name"`
		License struct {
			Text string `toml:"text"`
			File string `toml:"file"`
		} `toml:"license"`
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

	// Try to find pyproject.toml in the same directory to get project name and license
	dir := filepath.Dir(filePath)
	pyprojectPath := filepath.Join(dir, "pyproject.toml")
	if data, err := os.ReadFile(pyprojectPath); err == nil {
		var config pyproject
		if err := toml.Unmarshal(data, &config); err == nil {
			if config.Project.Name != "" {
				projectName = config.Project.Name
			}
			if config.Project.License.Text != "" {
				projectLicense = NormalizeLicense(config.Project.License.Text)
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
