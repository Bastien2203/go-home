package main

import (
	"github.com/Bastien2203/go-home/shared/types"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type ValueConverter func(val any) any

// Define how gohome capability become homekit service/characteristic
type ServiceDef struct {
	Type           byte
	ServiceType    string
	CharType       string
	NewService     func() *service.S
	ValueConverter ValueConverter
}

var CapabilityRegistry = map[types.CapabilityType]ServiceDef{
	types.CapabilityTemperature: {
		Type:           accessory.TypeSensor,
		ServiceType:    service.TypeTemperatureSensor,
		CharType:       characteristic.TypeCurrentTemperature,
		NewService:     func() *service.S { return service.NewTemperatureSensor().S },
		ValueConverter: func(v any) any { return v },
	},
	types.CapabilityHumidity: {
		Type:           accessory.TypeSensor,
		ServiceType:    service.TypeHumiditySensor,
		CharType:       characteristic.TypeCurrentRelativeHumidity,
		NewService:     func() *service.S { return service.NewHumiditySensor().S },
		ValueConverter: func(v any) any { return v },
	},
	types.CapabilityBattery: {
		Type:           accessory.TypeSensor,
		ServiceType:    service.TypeBatteryService,
		CharType:       characteristic.TypeBatteryLevel,
		NewService:     func() *service.S { return service.NewBatteryService().S },
		ValueConverter: func(v any) any { return v },
	},
	types.CapabilityButtonEvent: {
		Type:        accessory.TypeProgrammableSwitch,
		ServiceType: service.TypeStatelessProgrammableSwitch,
		CharType:    characteristic.TypeProgrammableSwitchEvent,
		NewService:  func() *service.S { return service.NewStatelessProgrammableSwitch().S },
		ValueConverter: func(v any) any {
			s, ok := v.(string)
			if !ok {
				return characteristic.ProgrammableSwitchEventSinglePress
			}
			switch s {
			case "double_press":
				return characteristic.ProgrammableSwitchEventDoublePress
			case "long_press":
				return characteristic.ProgrammableSwitchEventLongPress
			default:
				return characteristic.ProgrammableSwitchEventSinglePress
			}
		},
	},
}
