package enricher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"siftyan/internal/parser"
	"strings"
	"sync"
	"time"
)

type pypiResponse struct {
	Info struct {
		License           string   `json:"license"`
		LicenseExpression string   `json:"license_expression"`
		Classifiers       []string `json:"classifiers"`
	} `json:"info"`
}

type PyPIEnricher struct {
	client *http.Client
	cache  sync.Map
	sem    chan struct{} // Semaphore to limit concurrent requests
}

func NewPyPIEnricher() *PyPIEnricher {
	return &PyPIEnricher{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		sem: make(chan struct{}, 5), // Limit to 5 concurrent requests
	}
}

// EnrichTree traverses the tree and enriches pip dependencies
func (e *PyPIEnricher) EnrichTree(root *parser.Dependency) {
	// Collect all pip deps that need enrichment
	var deps []*parser.Dependency
	collectPipDeps(root, &deps)

	seen := make(map[*parser.Dependency]bool)
	var wg sync.WaitGroup
	for _, dep := range deps {
		if !seen[dep] {
			seen[dep] = true
			wg.Add(1)
			// go routine for each dependency concurrently limited by semaphore
			go func(d *parser.Dependency) {
				defer wg.Done()

				e.sem <- struct{}{}        // acquire slot
				defer func() { <-e.sem }() // release slot

				license, err := e.Enrich(d.Name)
				if err == nil && license != "UNKNOWN" {
					d.License = parser.NormalizeLicense(license)
				}
			}(dep)
		}
	}
	wg.Wait()
}

func collectPipDeps(dep *parser.Dependency, out *[]*parser.Dependency) {
	if dep.Depth > 0 && dep.Ecosystem == "pip" &&
		(dep.License == "UNKNOWN" || dep.License == "") {
		*out = append(*out, dep)
	}
	for _, child := range dep.Dependencies {
		collectPipDeps(child, out)
	}
}

// Enrich fetches metadata for a package from PyPI
func (e *PyPIEnricher) Enrich(packageName string) (string, error) {
	// Check cache first
	if val, ok := e.cache.Load(packageName); ok {
		return val.(string), nil
	}

	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)

	resp, err := e.client.Get(url)
	if err != nil {
		e.cache.Store(packageName, "UNKNOWN")
		return "UNKNOWN", nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "UNKNOWN", nil
	}

	var data pypiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	var license string = "UNKNOWN"

	// For Pip, sometimes the "license" field is more reliable for simple names,
	// but classifiers provide standard categories.
	if data.Info.License != "" && data.Info.License != "UNKNOWN" && data.Info.License != "null" {
		license = data.Info.License
	} else if data.Info.LicenseExpression != "" && data.Info.LicenseExpression != "null" {
		license = data.Info.LicenseExpression
	} else {
		// Scan classifiers for license info
		for _, c := range data.Info.Classifiers {
			// License :: OSI Approved :: BSD License e.g.
			if strings.HasPrefix(c, "License ::") {
				parts := strings.Split(c, " :: ")
				if len(parts) >= 3 {
					license = parts[len(parts)-1]
					break
				}
			}
		}
	}

	// Store in cache
	e.cache.Store(packageName, license)

	return license, nil
}
