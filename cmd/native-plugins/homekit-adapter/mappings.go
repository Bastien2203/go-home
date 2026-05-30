package main

import (
	"github.com/Bastien2203/go-home/shared/types"
	"github.com/absmach/senml"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type ValueConverter func(*senml.Record) any

// Define how gohome capability become homekit service/characteristic
type ServiceDef struct {
	Type           byte
	ServiceType    string
	CharType       string
	NewService     func() *service.S
	ValueConverter ValueConverter
}

var RecordRegistry = map[types.HubRecordName]ServiceDef{
	types.HubRecordTemperature: {
		Type:        accessory.TypeSensor,
		ServiceType: service.TypeTemperatureSensor,
		CharType:    characteristic.TypeCurrentTemperature,
		NewService:  func() *service.S { return service.NewTemperatureSensor().S },
		ValueConverter: func(r *senml.Record) any {
			if r.Value == nil {
				return nil
			}
			return *r.Value
		},
	},
	types.HubRecordHumidity: {
		Type:        accessory.TypeSensor,
		ServiceType: service.TypeHumiditySensor,
		CharType:    characteristic.TypeCurrentRelativeHumidity,
		NewService:  func() *service.S { return service.NewHumiditySensor().S },
		ValueConverter: func(r *senml.Record) any {
			if r.Value == nil {
				return nil
			}
			return *r.Value
		},
	},
	types.HubRecordBattery: {
		Type:        accessory.TypeSensor,
		ServiceType: service.TypeBatteryService,
		CharType:    characteristic.TypeBatteryLevel,
		NewService:  func() *service.S { return service.NewBatteryService().S },
		ValueConverter: func(r *senml.Record) any {
			if r.Value == nil {
				return nil
			}
			return *r.Value
		},
	},
	types.HubRecordButtonEvent: {
		Type:        accessory.TypeProgrammableSwitch,
		ServiceType: service.TypeStatelessProgrammableSwitch,
		CharType:    characteristic.TypeProgrammableSwitchEvent,
		NewService:  func() *service.S { return service.NewStatelessProgrammableSwitch().S },
		ValueConverter: func(r *senml.Record) any {
			if r.StringValue == nil {
				return nil
			}
			switch types.HubEvent(*r.StringValue) {
			case types.HubEventButtonPress:
				return characteristic.ProgrammableSwitchEventSinglePress
			case types.HubEventButtonDoublePress:
				return characteristic.ProgrammableSwitchEventDoublePress
			case types.HubEventButtonLongPress:
				return characteristic.ProgrammableSwitchEventLongPress
			default:
				return characteristic.ProgrammableSwitchEventSinglePress
			}
		},
	},
}
