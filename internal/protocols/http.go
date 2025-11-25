package protocols

import (
	"encoding/json"
	"fmt"
	"gohome/shared/types"
)

type HttpParser struct{}

func NewHttpParser() *HttpParser {
	return &HttpParser{}
}

func (d *HttpParser) AddressType() types.AddressType {
	return types.BasicAddress
}

func (s *HttpParser) ID() string {
	return "http"
}

func (d *HttpParser) Name() string {
	return "Http"
}

func (d *HttpParser) Parse(rawData []byte) ([]*types.Capability, error) {
	var capabilities []*types.Capability
	if err := json.Unmarshal(rawData, &capabilities); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to capabilities: %w", err)
	}
	return capabilities, nil
}
