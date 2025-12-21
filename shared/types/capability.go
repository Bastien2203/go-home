package types

type Capability struct {
	Name  CapabilityType `json:"name"`
	Value any            `json:"value"`
	Type  ValueType      `json:"type"` // ex: "float", "bool", "string"
	Unit  Unit           `json:"unit,omitempty"`
}

type CapabilityType string

const (
	CapabilityTemperature CapabilityType = "temperature"
	CapabilityHumidity    CapabilityType = "humidity"
	CapabilityBattery     CapabilityType = "battery_level"
	CapabilityButtonEvent CapabilityType = "button_event"
)

type ValueType string

const (
	TypeFloat  ValueType = "float"
	TypeInt    ValueType = "int"
	TypeBool   ValueType = "bool"
	TypeString ValueType = "string"
	TypeBytes  ValueType = "bytes"
)
