package engine

import (
	"siftyan/internal/parser"
)

// ConflictType represents the type of license conflict detected
type ConflictType string

const (
	CopyleftPropagation ConflictType = "Copyleft Propagation"
	NetworkCopyleft     ConflictType = "Network Copyleft"
	LicenseAmbiguity    ConflictType = "License Ambiguity"
)

type Conflict struct {
	Type        ConflictType
	Path        []string
	Description string
	Impact      string // HIGH, MEDIUM, LOW
	Suggestions []string
}

type ConflictDetector struct {
	registry *SPDXRegistry
}

func NewConflictDetector() *ConflictDetector {
	return &ConflictDetector{
		registry: GetSPDXRegistry(),
	}
}

// Detect analyzes the dependency tree for conflicts
func (d *ConflictDetector) Detect(root *parser.Dependency, model string) []Conflict {
	conflicts := []Conflict{}
	d.traverse(root, []string{root.Name}, &conflicts, model)
	return conflicts
}

func (d *ConflictDetector) traverse(dep *parser.Dependency, currentPath []string, conflicts *[]Conflict, model string) {
	// Detect Network Copyleft (AGPL) - affects even SaaS
	if d.registry.GetType(dep.License) == NetworkCopyleftLT {
		*conflicts = append(*conflicts, Conflict{
			Type:        NetworkCopyleft,
			Path:        append([]string{}, currentPath...),
			Description: "AGPL-3.0 license detected. This may require you to share your source code even if used as a network service.",
			Impact:      "HIGH",
			Suggestions: []string{"Replace this dependency", "Consult with legal team"},
		})
	}

	// Detect Copyleft Propagation (GPL in Permissive)
	// TODO: Check the root project's license
	// For simplicity, we assume the root is Permissive if not otherwise specified
	if d.registry.GetType(dep.License) == StrongCopyleftLT {
		*conflicts = append(*conflicts, Conflict{
			Type:        CopyleftPropagation,
			Path:        append([]string{}, currentPath...),
			Description: "GPL license detected in a project that might be intended to be permissive.",
			Impact:      "HIGH",
			Suggestions: []string{"Replace this dependency", "Change your project license to GPL"},
		})
	}

	// Detect License Ambiguity
	if dep.License == "UNKNOWN" || dep.License == "" {
		*conflicts = append(*conflicts, Conflict{
			Type:        LicenseAmbiguity,
			Path:        append([]string{}, currentPath...),
			Description: "No license detected for this dependency.",
			Impact:      "MEDIUM",
			Suggestions: []string{"Manually verify the license", "Add license info to lockfile if possible"},
		})
	}

	for _, child := range dep.Dependencies {
		newPath := append(currentPath, child.Name)
		d.traverse(child, newPath, conflicts, model)
	}
}
