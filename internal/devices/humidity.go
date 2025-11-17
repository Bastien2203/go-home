package devices

import (
	"gohome/internal/core"
	"time"
)

type HumiditySensor struct {
	*core.BaseDevice
}

func NewHumiditySensor(id string, name string, parser core.Parser) *HumiditySensor {
	return &HumiditySensor{BaseDevice: core.NewBaseDevice(id, name, core.HumiditySensorType, parser)}
}

func (s *HumiditySensor) PushHumidity(t float64) error {
	s.SetState("humidity", t)
	s.SetState("last_update", time.Now().UTC())
	return nil
}

func (s *HumiditySensor) GetHumidity() (float64, time.Time, bool) {
	st := s.State()
	v, ok := st["humidity"].(float64)
	if !ok {
		return 0, time.Time{}, false
	}
	lu, _ := st["last_update"].(time.Time)
	return v, lu, true
}
