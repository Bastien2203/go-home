package protocols

import (
	"errors"
	"fmt"
	"time"

	"github.com/Bastien2203/go-home/shared/types"
	"github.com/absmach/senml"
)

type SwitchBotParser struct {
	lastPayloads map[string]string
	timestamp    time.Time
}

func NewSwitchBotParser() *SwitchBotParser {
	return &SwitchBotParser{
		lastPayloads: make(map[string]string),
	}
}

func (p *SwitchBotParser) Name() string {
	return "switchbot"
}

func (p *SwitchBotParser) CanParse() bool {
	return true
}

func (d *SwitchBotParser) ClearCache() {
	if time.Now().After(d.timestamp.Add(TTL)) {
		d.lastPayloads = make(map[string]string)
		d.timestamp = time.Now()
	}
}

func (p *SwitchBotParser) Parse(address string, payload []byte) ([]senml.Record, bool, error) {
	p.ClearCache()
	encrypted := (payload[0] & 0b10000000) != 0

	if encrypted {
		return nil, false, fmt.Errorf("encrypted switchbot payload not supported for now")
	}

	payloadStr := string(payload)
	if last, ok := p.lastPayloads[address]; ok {
		if last == payloadStr {
			return nil, true, nil
		}
	}
	p.lastPayloads[address] = payloadStr

	modelChar := payload[0] & 0x7F

	var records []senml.Record
	var err error

	switch modelChar {
	case ModelMeter, ModelMeterPlus:
		records, err = parseMeter(payload)
		if err != nil {
			return nil, false, err
		}
	case ModelCurtain:
		return nil, false, fmt.Errorf("doesnt support switchbot curtain for now")
	case ModelMotionSensor:
		return nil, false, fmt.Errorf("doesnt support switchbot motion sensor for now")
	case ModelContactSensor:
		return nil, false, fmt.Errorf("doesnt support switchbot contact sensor for now")

	default:
		return nil, false, fmt.Errorf("doesnt support switchbot model: %c", modelChar)
	}

	return records, false, nil
}

func parseMeter(data []byte) ([]senml.Record, error) {
	if len(data) < 6 {
		return nil, errors.New("invalid meter data length")
	}

	records := []senml.Record{}

	// 1. Battery
	battVal := float64(data[2] & 0x7F)
	records = append(records, senml.Record{
		Name:  string(types.HubRecordBattery),
		Unit:  string(types.Percentage),
		Value: &battVal,
	})

	// 2. Temperature
	tempFrac := float64(data[3]&0x0F) / 10.0
	tempInt := float64(data[4] & 0x7F)
	tempVal := tempInt + tempFrac
	if (data[4] & 0x80) == 0 {
		tempVal = -tempVal
	}
	records = append(records, senml.Record{
		Name:  types.HubRecordTemperature.ToString(),
		Unit:  types.CelsiusDegree.ToString(),
		Value: &tempVal,
	})

	// 3. Humidity
	humVal := float64(data[5] & 0x7F)
	records = append(records, senml.Record{
		Name:  string(types.HubRecordHumidity),
		Unit:  string(types.HumidityPercentage),
		Value: &humVal,
	})

	return records, nil
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
