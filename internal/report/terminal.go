package report

import (
	"fmt"
	"siftyan/internal/engine"
	"strings"
)

// TerminalRenderer prints conflicts in plain English
type TerminalRenderer struct{}

func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{}
}

func (r *TerminalRenderer) Render(conflicts []engine.Conflict) {
	if len(conflicts) == 0 {
		fmt.Println("OK: No license conflicts detected!")
		return
	}

	fmt.Printf("WARNING:  Found %d license conflicts\n\n", len(conflicts))

	for i, c := range conflicts {
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
