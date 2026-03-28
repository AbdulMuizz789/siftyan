package enricher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
}

func NewPyPIEnricher() *PyPIEnricher {
	return &PyPIEnricher{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// Enrich fetches metadata for a package from PyPI
func (e *PyPIEnricher) Enrich(packageName string) (string, error) {
	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)

	resp, err := e.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "UNKNOWN", nil
	}

	var data pypiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	// For Pip, sometimes the "license" field is more reliable for simple names,
	// but classifiers provide standard categories.
	if data.Info.License != "" && data.Info.License != "UNKNOWN" && data.Info.License != "null" {
		return data.Info.License, nil
	}

	if data.Info.LicenseExpression != "" && data.Info.LicenseExpression != "null" {
		return data.Info.LicenseExpression, nil
	}

	// Scan classifiers for license info
	for _, c := range data.Info.Classifiers {
		// License :: OSI Approved :: BSD License e.g.
		if strings.HasPrefix(c, "License ::") {
			parts := strings.Split(c, " :: ")
			if len(parts) >= 3 {
				return parts[len(parts)-1], nil
			}
		}
	}

	return "UNKNOWN", nil
}
