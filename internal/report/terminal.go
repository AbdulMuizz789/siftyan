package report

import (
	"fmt"
	"siftyan/internal/engine"
	"strings"
)

// TerminalRenderer prints conflicts in plain English
type TerminalRenderer struct {
	conflicts []engine.Conflict
}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{
		conflicts: make([]engine.Conflict, 0),
	}
}

func (r *TerminalRenderer) OnConflictFound(c engine.Conflict) {
	r.conflicts = append(r.conflicts, c)
}

func (r *TerminalRenderer) Render() {
	if len(r.conflicts) == 0 {
		fmt.Println("OK: No license conflicts detected!")
		return
	}

	fmt.Printf("WARNING:  Found %d license conflicts\n\n", len(r.conflicts))

	for i, c := range r.conflicts {
		fmt.Printf("CONFLICT %d — %s\n", i+1, c.Type)
		fmt.Printf("Path: %s\n", strings.Join(c.Path, " → "))
		fmt.Printf("What this means: %s\n", c.Description)
		fmt.Printf("Impact: %s\n", c.Impact)
		fmt.Printf("Suggested actions:\n")
		for _, s := range c.Suggestions {
			fmt.Printf("  - %s\n", s)
		}
		fmt.Println()
	}
}
