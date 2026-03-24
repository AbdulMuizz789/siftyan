package parser

import (
	"os"
)

// LockfileParser is the strategy interface for parsing different lockfiles
type LockfileParser interface {
	Parse(filePath string) (*Dependency, error)
}

// BaseParser provides common functionality for all concrete parsers
type BaseParser struct {
	// Concrete parsers will implement their own Decode logic
	Decoder func(data []byte) (*Dependency, error)
}

func (p *BaseParser) ParseWith(filePath string) (*Dependency, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	dep, err := p.Decoder(data)
	if err != nil {
		return nil, err
	}

	return dep, nil
}
