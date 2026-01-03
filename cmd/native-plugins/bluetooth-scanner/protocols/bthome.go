package protocols

import (
	"fmt"
	"time"

	"log"

	"github.com/Bastien2203/bthomev2"
	bthomev2_types "github.com/Bastien2203/bthomev2/types"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

var BthomeUUID = bluetooth.New16BitUUID(uint16(bthomev2_types.ServiceDataUUID))

var TTL = 1 * time.Hour

type BthomeParser struct {
	cache     map[string]uint8
	timestamp time.Time
}

func NewBthomeParser() *BthomeParser {
	return &BthomeParser{
		timestamp: time.Now(),
		cache:     make(map[string]uint8),
	}
}

func (d *BthomeParser) Name() string {
	return "bthome"
}

func (d *BthomeParser) CanParse() bool {
	return true
}

func (d *BthomeParser) ClearCache() {
	if time.Now().After(d.timestamp.Add(TTL)) {
		d.cache = make(map[string]uint8)
		d.timestamp = time.Now()
	}
}

// Returns list of capabilities, boolean false if packet is duplicated, and error
func (d *BthomeParser) Parse(address string, payload []byte) ([]*types.Capability, bool, error) {
	d.ClearCache()
	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, false, nil
	}

	pid, found := data[bthomev2_types.PacketID]
	if found {
		pidValue, ok := pid.Float64()
		if !ok {
			return nil, false, fmt.Errorf("invalid packet id value")
		}
		// If packet already seen, ignore it
		if entry, ok := d.cache[address]; ok {
			if entry == uint8(pidValue) {
				return nil, true, nil
			}
		}
		// Else, update cache
		d.cache[address] = uint8(pidValue)
	}

	capabilities := make([]*types.Capability, 0)
	for _, m := range data {
		if c := CreateCapability(m); c != nil {
			capabilities = append(capabilities, c)
		}
	}

	return capabilities, false, nil
}

func UnitFromBthome(u bthomev2_types.Unit) types.Unit {
	switch u {
	case bthomev2_types.CelsiusDegree:
		return types.UnitCelsius
	case bthomev2_types.Percentage:
		return types.UnitPercent
	case bthomev2_types.Volt:
		return types.UnitVolt
	default:
		return types.NoUnit
	}
}

func CreateCapability(m bthomev2_types.Measurement) *types.Capability {
	name, ok := PropertyToCapability[m.Property]
	if !ok {
		return nil
	}

	var value any
	var t types.ValueType
	switch v := m.Value.(type) {
	case bthomev2_types.NumberValue:
		value = v.Number
		t = types.TypeFloat
	case bthomev2_types.BinaryValue:
		value = v.Boolean
		t = types.TypeBool

	case bthomev2_types.TextValue:
		value = v.Text
		t = types.TypeString
	case bthomev2_types.RawValue:
		value = v.Raw
		t = types.TypeBytes
	case bthomev2_types.EventValue:
		// TODO: add type = types.TypeEvent when core will be updated
		value = v.Event
	}

	return &types.Capability{
		Name:  name,
		Value: value,
		Type:  t,
		Unit:  UnitFromBthome(m.Unit),
	}
}

var PropertyToCapability = map[bthomev2_types.Property]types.CapabilityType{
	// bthomev2_types.PacketID: nil,
	// Sensors
	// bthomev2_types.SensorAcceleration:    nil,
	bthomev2_types.SensorBattery: types.CapabilityBattery,
	// bthomev2_types.SensorChannel:         nil,
	// bthomev2_types.SensorCO2:             nil,
	// bthomev2_types.SensorConductivity:    nil,
	// bthomev2_types.SensorCount:           nil,
	// bthomev2_types.SensorCurrent:         nil,
	// bthomev2_types.SensorDewPoint:        nil,
	// bthomev2_types.SensorDirection:       nil,
	// bthomev2_types.SensorDistanceMM:      nil,
	// bthomev2_types.SensorDistanceM:       nil,
	// bthomev2_types.SensorDuration:        nil,
	// bthomev2_types.SensorEnergy:          nil,
	// bthomev2_types.SensorGas:             nil,
	// bthomev2_types.SensorGyroscope:       nil,
	bthomev2_types.SensorHumidity: types.CapabilityHumidity,
	// bthomev2_types.SensorIlluminance:     nil,
	// bthomev2_types.SensorMassKG:          nil,
	// bthomev2_types.SensorMassLB:          nil,
	// bthomev2_types.SensorMoisture:        nil,
	// bthomev2_types.SensorPM2_5:           nil,
	// bthomev2_types.SensorPM10:            nil,
	// bthomev2_types.SensorPower:           nil,
	// bthomev2_types.SensorPrecipitation:   nil,
	// bthomev2_types.SensorPressure:        nil,
	// bthomev2_types.SensorRaw:             nil,
	// bthomev2_types.SensorRotation:        nil,
	// bthomev2_types.SensorRotational:      nil,
	// bthomev2_types.SensorSpeed:           nil,
	bthomev2_types.SensorTemperature: types.CapabilityTemperature,
	// bthomev2_types.SensorText:            nil,
	// bthomev2_types.SensorTimestamp:       nil,
	// bthomev2_types.SensorTVOC:            nil,
	// bthomev2_types.SensorVoltage:         nil,
	// bthomev2_types.SensorVolume:          nil,
	// bthomev2_types.SensorVolumeML:        nil,
	// bthomev2_types.SensorVolumeStorage:   nil,
	// bthomev2_types.SensorVolumeFlow:      nil,
	// bthomev2_types.SensorUV:              nil,
	// bthomev2_types.SensorWater:           nil,
	// bthomev2_types.SensorBatteryCharging: nil,
	// bthomev2_types.SensorCarbonMonoxide:  nil,
	// bthomev2_types.SensorCold:            nil,
	// bthomev2_types.SensorConnectivity:    nil,
	// bthomev2_types.SensorDoor:            nil,
	// bthomev2_types.SensorGarageDoor:      nil,
	// bthomev2_types.SensorGenericBoolean:  nil,
	// bthomev2_types.SensorHeat:            nil,
	// bthomev2_types.SensorLight:           nil,
	// bthomev2_types.SensorLock:            nil,
	// bthomev2_types.SensorMotion:          nil,
	// bthomev2_types.SensorMoving:          nil,
	// bthomev2_types.SensorOccupancy:       nil,
	// bthomev2_types.SensorOpening:         nil,
	// bthomev2_types.SensorPlug:            nil,
	// bthomev2_types.SensorPresence:        nil,
	// bthomev2_types.SensorProblem:         nil,
	// bthomev2_types.SensorRunning:         nil,
	// bthomev2_types.SensorSafety:          nil,
	// bthomev2_types.SensorSmoke:           nil,
	// bthomev2_types.SensorSound:           nil,
	// bthomev2_types.SensorTamper:          nil,
	// bthomev2_types.SensorVibration:       nil,
	// bthomev2_types.SensorWindow:          nil,
	bthomev2_types.SensorButtonEvent: types.CapabilityButtonEvent,
	// bthomev2_types.SensorDimmerEvent:     nil,
}
