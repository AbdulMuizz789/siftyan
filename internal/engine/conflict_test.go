package engine

import (
	"siftyan/internal/parser"
	"testing"
)

func TestConflictDetector(t *testing.T) {
	root, _ := parser.NewDependencyBuilder().
		Name("test-root").
		License("MIT").
		Ecosystem("npm").
		Build()

	libA, _ := parser.NewDependencyBuilder().
		Name("libA").
		License("GPL-3.0").
		Ecosystem("npm").
		Build()

	libB, _ := parser.NewDependencyBuilder().
		Name("libB").
		License("AGPL-3.0").
		Ecosystem("npm").
		Build()

	root.Dependencies = append(root.Dependencies, libA, libB)

	detector := NewConflictDetector()
	conflicts := detector.Detect(root, "binary")

	if len(conflicts) != 2 {
		t.Errorf("Expected 2 conflicts, got %d", len(conflicts))
	}

	if conflicts[0].Type != CopyleftPropagation {
		t.Errorf("Expected first conflict to be Copyleft Propagation, got %s", conflicts[0].Type)
	}

	if conflicts[1].Type != NetworkCopyleft {
		t.Errorf("Expected second conflict to be Network Copyleft, got %s", conflicts[1].Type)
	}
}
