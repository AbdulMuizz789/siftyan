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
	LinkingException    ConflictType = "Copyleft with Linking Exception"
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
	rules := GetRulesForModel(d.model)
	d.traverse(root, []string{root.Name}, rules)
}

func (d *ConflictDetector) traverse(dep *parser.Dependency, currentPath []string, rules []Rule) {
	licenseType := d.registry.GetType(dep.License)

	// Apply model-specific rules
	// TODO: Check the root project's license
	// For simplicity, we assume the root is Permissive if not otherwise specified
	for _, rule := range rules {
		if licenseType == rule.TargetType {
			ctype := CopyleftPropagation
			switch licenseType {
			case NetworkCopyleftLT:
				ctype = NetworkCopyleft
			case WeakCopyleftLT:
				ctype = LinkingException
			}

			d.notify(Conflict{
				Type:        ctype,
				Path:        append([]string{}, currentPath...),
				Description: rule.Description,
				Impact:      rule.Impact,
				Suggestions: rule.Suggestions,
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
		d.traverse(child, newPath, rules)
	}
}
