package report

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"siftyan/internal/engine"
	"siftyan/internal/parser"
	"text/template"
)

//go:embed report.html
var reportTemplate string

type HTMLRenderer struct {
	Tree      *parser.Dependency
	Conflicts []engine.Conflict
}

func NewHTMLRenderer(tree *parser.Dependency) *HTMLRenderer {
	return &HTMLRenderer{
		Tree:      tree,
		Conflicts: make([]engine.Conflict, 0),
	}
}

func (r *HTMLRenderer) OnConflictFound(c engine.Conflict) {
	r.Conflicts = append(r.Conflicts, c)
}

func (r *HTMLRenderer) WriteReport(outputPath string) error {
	tmpl, err := template.New("report").Parse(reportTemplate)
	if err != nil {
		return err
	}

	treeJSON, err := json.Marshal(r.Tree)
	if err != nil {
		return fmt.Errorf("failed to marshal dependency tree: %w", err)
	}
	conflictsJSON, err := json.Marshal(r.Conflicts)
	if err != nil {
		return fmt.Errorf("failed to marshal conflicts: %w", err)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data := struct {
		Conflicts     []engine.Conflict
		TreeJSON      string
		ConflictsJSON string
	}{
		Conflicts:     r.Conflicts,
		TreeJSON:      string(treeJSON),
		ConflictsJSON: string(conflictsJSON),
	}

	return tmpl.Execute(f, data)
}
