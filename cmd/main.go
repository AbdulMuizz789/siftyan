package main

import (
	"fmt"
	"os"
	"siftyan/internal/engine"
	"siftyan/internal/parser"
	"siftyan/internal/report"

	"github.com/spf13/cobra"
)

var model string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "siftyan",
		Short: "Siftyan - A local-first license conflict detector",
		Long:  `Siftyan scans your project's dependency tree and detects license conflicts.`,
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan for license conflicts",
		Run: func(cmd *cobra.Command, args []string) {
			dir, _ := os.Getwd()
			files, err := parser.Detect(dir)
			if err != nil {
				fmt.Printf("Error detecting lockfiles: %v\n", err)
				return
			}

			if len(files) == 0 {
				fmt.Println("No supported lockfiles found (package-lock.json or requirements.txt)")
				return
			}

			allConflicts := []engine.Conflict{}
			detector := engine.NewConflictDetector()
			renderer := report.NewTerminalRenderer()

			for _, file := range files {
				fmt.Printf("Scanning %s...\n", file)
				p, err := parser.NewForFile(file)
				if err != nil {
					fmt.Printf("Error creating parser: %v\n", err)
					continue
				}

				root, err := p.Parse(file)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", file, err)
					continue
				}

				conflicts := detector.Detect(root, model)
				allConflicts = append(allConflicts, conflicts...)
			}

			renderer.Render(allConflicts)
		},
	}

	scanCmd.Flags().StringVarP(&model, "model", "m", "internal", "Distribution model (saas|binary|internal)")
	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
