package parser

import (
	"strings"
)

// NormalizeLicense canonicalizes license strings into SPDX identifiers
// TODO: This is a very basic cleanup. It will be expanded later
func NormalizeLicense(license string) string {
	license = strings.TrimSpace(license)
	if license == "" {
		return "UNKNOWN"
	}

	// Basic normalization: uppercase for common licenses
	l := strings.ToUpper(license)
	if strings.Contains(l, "MIT") {
		return "MIT"
	}
	if strings.Contains(l, "APACHE") {
		return "Apache-2.0"
	}
	if strings.Contains(l, "BSD") {
		if strings.Contains(l, "3") {
			return "BSD-3-Clause"
		}
		if strings.Contains(l, "2") {
			return "BSD-2-Clause"
		}
		return "BSD-3-Clause" // Default to BSD-3
	}
	if strings.Contains(l, "GPL") {
		if strings.Contains(l, "3") {
			if strings.Contains(l, "AFFERO") || strings.Contains(l, "AGPL") {
				return "AGPL-3.0"
			}
			if strings.Contains(l, "LOWER") || strings.Contains(l, "LGPL") {
				return "LGPL-3.0"
			}
			return "GPL-3.0"
		}
		if strings.Contains(l, "2") {
			if strings.Contains(l, "LOWER") || strings.Contains(l, "LGPL") {
				return "LGPL-2.1"
			}
			return "GPL-2.0"
		}
		return "GPL-3.0"
	}
	if strings.Contains(l, "MPL") {
		return "MPL-2.0"
	}

	return license
}
