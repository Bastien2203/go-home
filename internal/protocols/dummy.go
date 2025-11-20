package protocols

import (
	"encoding/json"
	"fmt"

	"gohome/internal/core"
)

// DummyParser implémente le protocol pour les données de test
type DummyParser struct{}

func NewDummyParser() *DummyParser {
	return &DummyParser{}
}

func (d *DummyParser) AddressType() core.AddressType {
	return core.BLEAddress
}

func (s *DummyParser) ID() string {
	return "dummy"
}

func (d *DummyParser) Name() string {
	return "Dummy"
}

func (d *DummyParser) Parse(rawData any) ([]*core.Capability, error) {
	var rawMap map[string]any
	if err := json.Unmarshal([]byte(rawData.(string)), &rawMap); err != nil {
		return nil, fmt.Errorf("failed to parse raw data: %v", err)
	}

	deviceType := core.CapabilityTemperature
	if dt, exists := rawMap["type"]; exists {
		deviceType = core.CapabilityType(dt.(string))
	}

	value := rawMap["value"]
	if value == nil {
		value = 0.0
	}

	valueType := core.TypeFloat
	switch value.(type) {
	case float64:
		valueType = core.TypeFloat
	case int:
		valueType = core.TypeInt
	case bool:
		valueType = core.TypeBool
	case string:
		valueType = core.TypeString
	}

	unit := core.UnitCelsius
	if u, exists := rawMap["unit"]; exists {
		unit = core.Unit(u.(string))
	}

	return []*core.Capability{
		{
			Name:  deviceType,
			Value: value,
			Type:  valueType,
			Unit:  unit,
		},
	}, nil
}
