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

// ConflictObserver defines the interface for receiving conflict events
type ConflictObserver interface {
	OnConflictFound(c Conflict)
}

// Option is a functional option for ConflictDetector
type Option func(*ConflictDetector)

type ConflictDetector struct {
	registry  *SPDXRegistry
	observers []ConflictObserver
	model     string
}

func NewConflictDetector(opts ...Option) *ConflictDetector {
	d := &ConflictDetector{
		registry: GetSPDXRegistry(),
		model:    "internal", // default
	}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Functional Options
func WithModel(model string) Option {
	return func(d *ConflictDetector) {
		d.model = model
	}
}

func WithObserver(o ConflictObserver) Option {
	return func(d *ConflictDetector) {
		d.observers = append(d.observers, o)
	}
}

func (d *ConflictDetector) notify(c Conflict) {
	for _, o := range d.observers {
		o.OnConflictFound(c)
	}
}

// Detect analyzes the dependency tree for conflicts and notifies observers
func (d *ConflictDetector) Detect(root *parser.Dependency) {
	d.traverse(root, []string{root.Name})
}

func (d *ConflictDetector) traverse(dep *parser.Dependency, currentPath []string) {
	licenseType := d.registry.GetType(dep.License)

	// Detect Network Copyleft (AGPL) - Critical for SaaS and Binary
	if licenseType == NetworkCopyleftLT {
		d.notify(Conflict{
			Type:        NetworkCopyleft,
			Path:        append([]string{}, currentPath...),
			Description: "AGPL-3.0 detected. Requires source disclosure even for network services.",
			Impact:      "HIGH",
			Suggestions: []string{"Replace with a permissive alternative", "Isolate as a separate microservice (if legal allows)"},
		})
	}

	// Detect Copyleft Propagation - Depends on Model
	// TODO: Check the root project's license
	// For simplicity, we assume the root is Permissive if not otherwise specified
	if licenseType == StrongCopyleftLT {
		switch d.model {
		case "binary":
			// If model is binary, this is a HIGH impact conflict for permissive projects
			d.notify(Conflict{
				Type:        CopyleftPropagation,
				Path:        append([]string{}, currentPath...),
				Description: "Strong Copyleft detected in a binary distribution. This triggers the 'viral' clause.",
				Impact:      "HIGH",
				Suggestions: []string{"Replace this dependency", "Change your project license to GPL"},
			})
		case "saas":
			// In SaaS, standard GPL is often acceptable (unlike AGPL)
			d.notify(Conflict{
				Type:        CopyleftPropagation,
				Path:        append([]string{}, currentPath...),
				Description: "Strong Copyleft detected. Acceptable for internal SaaS use, but verify no client-side code is included.",
				Impact:      "LOW",
				Suggestions: []string{"Verify that this code is not shipped to the browser"},
			})
		}
	}

	// Detect License Ambiguity
	if dep.License == "UNKNOWN" || dep.License == "" {
		d.notify(Conflict{
			Type:        LicenseAmbiguity,
			Path:        append([]string{}, currentPath...),
			Description: "Missing license information. Legal risk cannot be automatically assessed.",
			Impact:      "MEDIUM",
			Suggestions: []string{"Manually check package repository", "Verify license in source headers", "Add license info to lockfile if possible"},
		})
	}

	for _, child := range dep.Dependencies {
		newPath := append(currentPath, child.Name)
		d.traverse(child, newPath)
	}
}
