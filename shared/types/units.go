package types

import (
	"github.com/Bastien2203/bthomev2"
)

type Unit string

const (
	UnitCelsius Unit = "celsius"
	UnitPercent Unit = "percent"
	UnitVolt    Unit = "volt"
	NoUnit      Unit = ""
)

func UnitFromBthome(u bthomev2.Unit) Unit {
	switch u {
	case bthomev2.CelsiusDegree:
		return UnitCelsius
	case bthomev2.Percentage:
		return UnitPercent
	case bthomev2.Volt:
		return UnitVolt
	default:
		return NoUnit
	}
}
