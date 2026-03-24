package parser

import (
	"os"
	"testing"
)

const sampleNpmLock = `{
  "name": "test-project",
  "version": "1.0.0",
  "lockfileVersion": 3,
  "packages": {
    "": {
      "name": "test-project",
      "version": "1.0.0",
      "license": "MIT"
    },
    "node_modules/lodash": {
      "version": "4.17.21",
      "license": "MIT"
    },
    "node_modules/express": {
      "version": "4.18.2",
      "license": "MIT"
    }
  }
}`

func TestNpmParser(t *testing.T) {
	tmpFile := "package-lock.json"
	err := os.WriteFile(tmpFile, []byte(sampleNpmLock), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary lockfile: %v", err)
	}
	defer os.Remove(tmpFile)

	parser := NewNpmParser()
	dep, err := parser.Parse(tmpFile)
	if err != nil {
		t.Fatalf("NpmParser.Parse failed: %v", err)
	}

	if dep.Name != "test-project" {
		t.Errorf("Expected root name 'test-project', got '%s'", dep.Name)
	}

	if dep.License != "MIT" {
		t.Errorf("Expected root license 'MIT', got '%s'", dep.License)
	}

	if len(dep.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(dep.Dependencies))
	}

	foundLodash := false
	for _, child := range dep.Dependencies {
		if child.Name == "lodash" {
			foundLodash = true
			if child.License != "MIT" {
				t.Errorf("Expected lodash license 'MIT', got '%s'", child.License)
			}
		}
	}

	if !foundLodash {
		t.Error("Did not find lodash in dependencies")
	}
}
