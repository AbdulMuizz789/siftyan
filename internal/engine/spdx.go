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
				"GPL-2.0":      StrongCopyleftLT,
				"GPL-3.0":      StrongCopyleftLT,
				"AGPL-3.0":     NetworkCopyleftLT,
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
