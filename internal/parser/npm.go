package parser

import (
	"encoding/json"
	"fmt"
	"strings"
)

// npmLockfile represents the structure of an npm lockfile (package-lock.json).
// See "https://docs.npmjs.com/cli/v8/configuring-npm/package-lock-json#file-format" for details.
type npmLockfile struct {
	Name            string                `json:"name"`
	Version         string                `json:"version"`
	LockfileVersion int                   `json:"lockfileVersion"`
	Packages        map[string]npmPackage `json:"packages"`
}

type npmPackage struct {
	// Version is the version of the package
	Version string `json:"version"`

	// Resolved is the URL or file path where the package was resolved from
	Resolved string `json:"resolved"`

	// Integrity is the integrity hash of the package content
	Integrity string `json:"integrity"`

	// License is the license information for the package
	License string `json:"license"`

	// Dependencies is a map of dependencies
	Dependencies map[string]string `json:"dependencies"`

	// DevDependencies is only present in the root package usually
	DevDependencies map[string]string `json:"devDependencies"`

	// Dev indicates if this package is a development dependency
	Dev bool `json:"dev"`
}

// NpmParser responsible for parsing npm lockfiles and extracting dependency information.
type NpmParser struct {
	BaseParser
	IncludeDev bool
}

func NewNpmParser() *NpmParser {
	p := &NpmParser{}
	p.Decoder = p.decode
	return p
}

func (p *NpmParser) WithIncludeDev(include bool) *NpmParser {
	p.IncludeDev = include
	return p
}

func (p *NpmParser) Parse(filePath string) (*Dependency, error) {
	return p.ParseWith(filePath)
}

func (p *NpmParser) decode(data []byte) (*Dependency, error) {
	var lock npmLockfile
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("failed to decode npm lockfile: %w", err)
	}

	// Build logical tree starting from root
	return p.buildRecursive(lock.Name, "", lock.Packages, 0, make(map[string]bool))
}

func (p *NpmParser) buildRecursive(name string, path string, packages map[string]npmPackage, depth int, visited map[string]bool) (*Dependency, error) {
	pkg, ok := packages[path]
	if !ok {
		return nil, fmt.Errorf("package not found at path: %s", path)
	}

	// Prevent infinite recursion in case of cycles (rare in physical layout but good safety)
	if visited[path] {
		// Return a simplified dependency to break cycle
		return NewDependencyBuilder().
			Name(name + " (cycle)").
			Version(pkg.Version).
			License(NormalizeLicense(pkg.License)).
			Ecosystem("npm").
			Depth(depth).
			Build()
	}
	visited[path] = true
	defer delete(visited, path)

	builder := NewDependencyBuilder().
		Name(name).
		Version(pkg.Version).
		License(NormalizeLicense(pkg.License)).
		Ecosystem("npm").
		Depth(depth)

	// Combine dependencies and devDependencies (if at root and includeDev is true)
	depsToResolve := make(map[string]string)
	for k, v := range pkg.Dependencies {
		depsToResolve[k] = v
	}
	if path == "" && p.IncludeDev {
		for k, v := range pkg.DevDependencies {
			depsToResolve[k] = v
		}
	}

	// Recursively add dependencies
	for depName := range depsToResolve {
		depPath := p.findPackagePath(path, depName, packages)
		if depPath != "" {
			childPkg := packages[depPath]

			// Respect IncludeDev flag for non-root dependencies
			if childPkg.Dev && !p.IncludeDev {
				continue
			}

			child, err := p.buildRecursive(depName, depPath, packages, depth+1, visited)
			if err == nil {
				builder.AddDependency(child)
			}
		}
	}

	return builder.Build()
}

func (p *NpmParser) findPackagePath(currentPath string, depName string, packages map[string]npmPackage) string {
	// npm looks for dependencies in node_modules/ at the current level, then parents
	parts := splitPath(currentPath)
	if currentPath == "" {
		parts = []string{}
	}

	// Check from deepest level to root
	for i := len(parts); i >= 0; i-- {
		prefix := ""
		if i > 0 {
			prefix = strings.Join(parts[:i], "/") + "/"
		}
		path := prefix + "node_modules/" + depName
		if _, ok := packages[path]; ok {
			return path
		}
	}

	return ""
}

func splitPath(path string) []string {
	// Handle both Windows and Unix paths if necessary, but npm uses "/"
	return strings.Split(path, "/")
}
