package protocols

import (
	"gohome/internal/core"
)

type HttpParser struct{}

func NewHttpParser() *HttpParser {
	return &HttpParser{}
}

func (d *HttpParser) AddressType() core.AddressType {
	return core.BasicAddress
}

func (s *HttpParser) ID() string {
	return "http"
}

func (d *HttpParser) Name() string {
	return "Http"
}

func (d *HttpParser) Parse(rawData any) ([]*core.Capability, error) {
	return rawData.([]*core.Capability), nil
}
