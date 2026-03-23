package parser

import (
	"errors"
	"strings"
)

// Dependency represents a single package in the dependency tree
type Dependency struct {
	// Name is the name of the package
	Name string `json:"name"`

	// Version is the version of the package
	Version string `json:"version"`

	// License is the license of the package discovered
	License string `json:"license"`

	// Ecosystem is the package manager ecosystem (e.g., npm, pip, maven, etc.)
	Ecosystem string `json:"ecosystem"`

	// Sub-dependencies of this package
	Dependencies []*Dependency `json:"dependencies,omitempty"`

	// Depth indicates how deep this dependency is in the tree (0 for root, 1 for direct dependencies, etc.)
	Depth int `json:"depth"`
}

// DependencyBuilder provides a fluent API for creating Dependency structs
type DependencyBuilder struct {
	dep *Dependency
	err error
}

func NewDependencyBuilder() *DependencyBuilder {
	return &DependencyBuilder{
		dep: &Dependency{
			Dependencies: make([]*Dependency, 0),
		},
	}
}

func (b *DependencyBuilder) Name(name string) *DependencyBuilder {
	if name == "" {
		b.err = errors.New("dependency name cannot be empty")
	}
	b.dep.Name = name
	return b
}

func (b *DependencyBuilder) Version(version string) *DependencyBuilder {
	b.dep.Version = version
	return b
}

func (b *DependencyBuilder) License(license string) *DependencyBuilder {
	b.dep.License = license
	return b
}

func (b *DependencyBuilder) Ecosystem(ecosystem string) *DependencyBuilder {
	b.dep.Ecosystem = strings.ToLower(ecosystem)
	return b
}

func (b *DependencyBuilder) Depth(depth int) *DependencyBuilder {
	b.dep.Depth = depth
	return b
}

func (b *DependencyBuilder) AddDependency(child *Dependency) *DependencyBuilder {
	if child != nil {
		b.dep.Dependencies = append(b.dep.Dependencies, child)
	}
	return b
}

func (b *DependencyBuilder) Build() (*Dependency, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.dep.Name == "" {
		return nil, errors.New("cannot build dependency: missing name")
	}
	return b.dep, nil
}
