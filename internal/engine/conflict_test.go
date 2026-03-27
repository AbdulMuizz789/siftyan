package engine

import (
	"siftyan/internal/parser"
	"testing"
)

type mockObserver struct {
	conflicts []Conflict
}

func (o *mockObserver) OnConflictFound(c Conflict) {
	o.conflicts = append(o.conflicts, c)
}

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

	obs := &mockObserver{}
	detector := NewConflictDetector(
		WithModel("binary"),
		WithObserver(obs),
	)
	detector.Detect(root)

	if len(obs.conflicts) != 2 {
		t.Errorf("Expected 2 conflicts, got %d", len(obs.conflicts))
	}

	foundGPL := false
	foundAGPL := false
	for _, c := range obs.conflicts {
		if c.Type == CopyleftPropagation {
			foundGPL = true
		}
		if c.Type == NetworkCopyleft {
			foundAGPL = true
		}
	}

	if !foundGPL {
		t.Error("Did not find Copyleft Propagation conflict")
	}
	if !foundAGPL {
		t.Error("Did not find Network Copyleft conflict")
	}
}

func TestModelLogic(t *testing.T) {
	root, _ := parser.NewDependencyBuilder().
		Name("test-root").
		License("GPL-3.0").
		Ecosystem("npm").
		Build()

	// Test SaaS model - GPL should be LOW impact or different description
	obsSaaS := &mockObserver{}
	detSaaS := NewConflictDetector(WithModel("saas"), WithObserver(obsSaaS))
	detSaaS.Detect(root)

	if len(obsSaaS.conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(obsSaaS.conflicts))
	}
	if obsSaaS.conflicts[0].Impact != "LOW" {
		t.Errorf("Expected LOW impact for GPL in SaaS, got %s", obsSaaS.conflicts[0].Impact)
	}

	// Test Binary model - GPL should be HIGH impact
	obsBin := &mockObserver{}
	detBin := NewConflictDetector(WithModel("binary"), WithObserver(obsBin))
	detBin.Detect(root)

	if len(obsBin.conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(obsBin.conflicts))
	}
	if obsBin.conflicts[0].Impact != "HIGH" {
		t.Errorf("Expected HIGH impact for GPL in Binary, got %s", obsBin.conflicts[0].Impact)
	}
}
