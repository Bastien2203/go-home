package types

type HubUnit string

func (u HubUnit) ToString() string {
	return string(u)
}

const (
	CelsiusDegree           HubUnit = "Cel"
	HumidityPercentage      HubUnit = "%RH"
	Percentage              HubUnit = "%"
	Volt                    HubUnit = "V"
	Ampere                  HubUnit = "A"
	Watt                    HubUnit = "W"
	KilowattHours           HubUnit = "kWh"
	MetersPerSecondSquared  HubUnit = "m/s2"
	PartsPerMillion         HubUnit = "ppm"
	Microsiemens            HubUnit = "S/m"
	Degree                  HubUnit = "deg"
	Millimeter              HubUnit = "mm"
	Meter                   HubUnit = "m"
	Second                  HubUnit = "s"
	CubicMeter              HubUnit = "m3"
	DegreesPerSecond        HubUnit = "deg/s"
	Lux                     HubUnit = "lx"
	Kilogram                HubUnit = "kg"
	Pound                   HubUnit = "lb"
	MicrogramsPerCubicMeter HubUnit = "ug/m3"
	Hectopascal             HubUnit = "hPa"
	RevolutionsPerMinute    HubUnit = "count/min"
	MetersPerSecond         HubUnit = "m/s"
	Liter                   HubUnit = "l"
	Milliliter              HubUnit = "ml"
	CubicMetersPerHour      HubUnit = "m3/h"
	Empty                   HubUnit = ""
)
