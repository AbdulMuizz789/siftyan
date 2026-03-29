package main

import (
	"fmt"
	"os"
	"siftyan/internal/engine"
	"siftyan/internal/enricher"
	"siftyan/internal/parser"
	"siftyan/internal/report"

	"github.com/spf13/cobra"
)

var model string
var reportPath string

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

			renderer := report.NewTerminalRenderer()
			pypiEnricher := enricher.NewPyPIEnricher()

			var htmlRenderer *report.HTMLRenderer
			var combinedRoot *parser.Dependency

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

				// Enrich pip dependencies
				pypiEnricher.EnrichTree(root)

				// Initialize combined root for HTML report
				if combinedRoot == nil {
					combinedRoot = root
					if reportPath != "" {
						htmlRenderer = report.NewHTMLRenderer(combinedRoot)
					}
				} else {
					combinedRoot.Dependencies = append(combinedRoot.Dependencies, root.Dependencies...)
				}

				// Create detector with observers
				opts := []engine.Option{
					engine.WithModel(model),
					engine.WithObserver(renderer),
				}
				if htmlRenderer != nil {
					opts = append(opts, engine.WithObserver(htmlRenderer))
				}

				detector := engine.NewConflictDetector(opts...)
				detector.Detect(root)
			}

			renderer.Render()

			if reportPath != "" && htmlRenderer != nil {
				if err := htmlRenderer.WriteReport(reportPath); err != nil {
					fmt.Printf("Error writing HTML report: %v\n", err)
				} else {
					fmt.Printf("HTML report generated at: %s\n", reportPath)
				}
			}
		},
	}

	scanCmd.Flags().StringVarP(&model, "model", "m", "internal", "Distribution model (saas|binary|internal)")
	scanCmd.Flags().StringVarP(&reportPath, "report", "r", "", "Path to generate HTML report (e.g., report.html)")
	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
