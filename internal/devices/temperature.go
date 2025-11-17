package devices

import (
	"gohome/internal/core"
	"time"
)

type TemperatureSensor struct {
	*core.BaseDevice
}

func NewTemperatureSensor(id string, name string, parser core.Parser) *TemperatureSensor {
	return &TemperatureSensor{BaseDevice: core.NewBaseDevice(id, name, core.TemperatureSensorType, parser)}
}

func (s *TemperatureSensor) PushTemperature(t float64) error {
	s.SetState("temperature", t)
	s.SetState("last_update", time.Now().UTC())
	return nil
}

func (s *TemperatureSensor) GetTemperature() (float64, time.Time, bool) {
	st := s.State()
	v, ok := st["temperature"].(float64)
	if !ok {
		return 0, time.Time{}, false
	}
	lu, _ := st["last_update"].(time.Time)
	return v, lu, true
}
