package protocols

import (
	"errors"
	"fmt"

	"github.com/Bastien2203/go-home/shared/types"
)

type SwitchBotParser struct {
}

func NewSwitchBotParser() *SwitchBotParser {
	return &SwitchBotParser{}
}

func (p *SwitchBotParser) Name() string {
	return "switchbot"
}

func (p *SwitchBotParser) CanParse() bool {
	return true
}

func (p *SwitchBotParser) Parse(address string, payload []byte) ([]*types.Capability, error) {
	encrypted := (payload[0] & 0b10000000) != 0

	if encrypted {
		return nil, fmt.Errorf("encrypted switchbot payload not supported for now")
	}
	modelChar := payload[0] & 0x7F

	switch modelChar {
	case ModelMeter, ModelMeterPlus:
		return parseMeter(payload)
	case ModelCurtain:
		return nil, fmt.Errorf("doesnt support switchbot curtain for now")
	case ModelMotionSensor:
		return nil, fmt.Errorf("doesnt support switchbot motion sensor for now")
	case ModelContactSensor:
		return nil, fmt.Errorf("doesnt support switchbot contact sensor for now")

	default:
		return nil, fmt.Errorf("doesnt support switchbot model: %c", modelChar)
	}
}

func parseMeter(data []byte) ([]*types.Capability, error) {
	if len(data) < 6 {
		return nil, errors.New("data meter invalide")
	}

	capabilites := make([]*types.Capability, 0, 3)

	capabilites = append(capabilites, &types.Capability{
		Name:  types.CapabilityBattery,
		Value: int(data[2] & 0x7F),
		Type:  types.TypeInt,
	})

	tempFrac := float64(data[3]&0x0F) / 10.0
	tempInt := float64(data[4] & 0x7F)
	temp := tempInt + tempFrac
	if (data[4] & 0x80) == 0 {
		temp = -temp
	}

	capabilites = append(capabilites, &types.Capability{
		Name:  types.CapabilityTemperature,
		Value: temp,
		Type:  types.TypeFloat,
		Unit:  types.UnitCelsius,
	})

	hum := int(data[5] & 0x7F)

	capabilites = append(capabilites, &types.Capability{
		Name:  types.CapabilityHumidity,
		Value: hum,
		Type:  types.TypeInt,
		Unit:  types.UnitPercent,
	})

	return capabilites, nil
}

const (
	ModelBot           = 'H' // WoHand
	ModelMeter         = 'T' // WoSensorTH
	ModelMeterPlus     = 'i' // WoSensorTH Plus
	ModelCurtain       = 'c' // WoCurtain
	ModelContactSensor = 'd' // WoContact
	ModelMotionSensor  = 's' // WoPresence
	ModelPlugMini      = 'g' // WoPlug
)
