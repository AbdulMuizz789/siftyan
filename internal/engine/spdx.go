package engine

import (
	"sync"
)

// LicenseType defines the category of a license
type LicenseType int

const (
	PermissiveLT LicenseType = iota
	WeakCopyleftLT
	StrongCopyleftLT
	NetworkCopyleftLT
	UnknownLT
)

// SPDXRegistry holds license information and compatibility rules
type SPDXRegistry struct {
	types map[string]LicenseType
}

var (
	instance *SPDXRegistry
	once     sync.Once
)

// GetSPDXRegistry returns the singleton instance of the registry
func GetSPDXRegistry() *SPDXRegistry {
	once.Do(func() {
		instance = &SPDXRegistry{
			types: map[string]LicenseType{
				"MIT":          PermissiveLT,
				"Apache-2.0":   PermissiveLT,
				"BSD-3-Clause": PermissiveLT,
				"BSD-2-Clause": PermissiveLT,
				"LGPL-3.0":     WeakCopyleftLT,
				"LGPL-2.1":     WeakCopyleftLT,
				"MPL-2.0":      WeakCopyleftLT,
				"GPL-2.0":      StrongCopyleftLT,
				"GPL-3.0":      StrongCopyleftLT,
				"AGPL-3.0":     NetworkCopyleftLT,
				"ISC":          PermissiveLT,
				"CC0-1.0":      PermissiveLT,
				"Unlicense":    PermissiveLT,
				"WTFPL":        PermissiveLT,
				"PSF-2.0":      PermissiveLT, // Python packages
				"Artistic-2.0": PermissiveLT,
				"EPL-2.0":      WeakCopyleftLT,
				"EUPL-1.2":     WeakCopyleftLT,
				"CDDL-1.0":     WeakCopyleftLT,
			},
		}
	})
	return instance
}

func (r *SPDXRegistry) GetType(license string) LicenseType {
	if t, ok := r.types[license]; ok {
		return t
	}
	return UnknownLT
}
