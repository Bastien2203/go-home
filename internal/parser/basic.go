package parser

import (
	"gohome/internal/core"
	"gohome/internal/scanner"
)

type BasicParser struct {
	scanner *scanner.HttpScanner
}

func NewBasicParser(scanner *scanner.HttpScanner) core.Parser {
	return &BasicParser{scanner: scanner}
}

func (p *BasicParser) Scanner() core.Scanner { return p.scanner }

func (p *BasicParser) Name() string { return "basic" }

func (p *BasicParser) Parse(raw core.Advertisment) (map[string]any, bool) {
	basicAdv, ok := raw.(*core.BasicAdvertisment)
	if !ok {
		return nil, false
	}

	return map[string]any{
		basicAdv.Type: basicAdv.Value,
	}, true
}
