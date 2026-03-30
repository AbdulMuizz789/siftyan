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

	// Root project info is in lock.Packages[""]
	rootPkg, ok := lock.Packages[""]
	if !ok {
		return nil, fmt.Errorf("root project info not found in npm lockfile packages")
	}

	rootBuilder := NewDependencyBuilder().
		Name(lock.Name).
		Version(rootPkg.Version).
		License(NormalizeLicense(rootPkg.License)).
		Ecosystem("npm").
		Depth(0)

	// Build the dependency tree
	// TODO: we only collect direct dependencies listed in packages.
	// A full tree builder would resolve transitive dependencies correctly
	for pkgPath, pkg := range lock.Packages {
		if pkgPath == "" {
			continue
		}

		// Only include direct dependencies for now (not nested in node_modules/...)
		// In npm v2/v3, "packages" keys look like "node_modules/name"
		if !isDirectDependency(pkgPath) {
			continue
		}

		// Filter out dev dependencies unless requested
		if pkg.Dev && !p.IncludeDev {
			continue
		}

		dep, err := NewDependencyBuilder().
			Name(cleanNpmPackageName(pkgPath)).
			Version(pkg.Version).
			License(NormalizeLicense(pkg.License)).
			Ecosystem("npm").
			Depth(1).
			Build()

		if err == nil {
			rootBuilder.AddDependency(dep)
		}
	}

	return rootBuilder.Build()
}

func isDirectDependency(path string) bool {
	// A direct dependency path looks like "node_modules/packageName"
	// Nested dependencies look like "node_modules/A/node_modules/B"
	// This is a simple heuristic.
	parts := splitPath(path)
	return len(parts) == 2 && parts[0] == "node_modules"
}

func cleanNpmPackageName(path string) string {
	parts := splitPath(path)
	if len(parts) >= 2 && parts[0] == "node_modules" {
		return parts[1]
	}
	return path
}

func splitPath(path string) []string {
	// Handle both Windows and Unix paths if necessary, but npm uses "/"
	return strings.Split(path, "/")
}
