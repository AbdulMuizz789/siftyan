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
	switch l {
	case "MIT":
		return "MIT"
	case "APACHE-2.0", "APACHE 2.0", "APACHE2.0":
		return "Apache-2.0"
	case "GPL-3.0", "GPL3.0", "GPLV3":
		return "GPL-3.0"
	case "GPL-2.0", "GPL2.0", "GPLV2":
		return "GPL-2.0"
	case "BSD-3-CLAUSE", "BSD 3-CLAUSE":
		return "BSD-3-Clause"
	case "BSD-2-CLAUSE", "BSD 2-CLAUSE":
		return "BSD-2-Clause"
	case "AGPL-3.0", "AGPL3.0":
		return "AGPL-3.0"
	case "LGPL-3.0", "LGPL3.0":
		return "LGPL-3.0"
	}

	return license
}
