package protocols

import (
	senmltypes "bluetooth-scanner/senml_types"
	"fmt"
	"time"

	"log"

	"github.com/Bastien2203/bthomev2"
	bthomev2_types "github.com/Bastien2203/bthomev2/types"
	"github.com/absmach/senml"

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
func (d *BthomeParser) Parse(address string, payload []byte) ([]senml.Record, bool, error) {
	d.ClearCache()
	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, false, err
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

	records := []senml.Record{}

	for _, m := range data {
		if m.Property == bthomev2_types.PacketID {
			continue
		}
		if record := CreateSenMLRecord(m); record != nil {
			records = append(records, *record)
		}
	}

	return records, false, nil
}

func CreateSenMLRecord(m bthomev2_types.Measurement) *senml.Record {
	name, ok := propertyToName[m.Property]
	if !ok {
		return nil
	}

	rec := &senml.Record{
		Name: name.ToString(),
		Unit: UnitToSenML(m.Property, m.Unit).ToString(),
	}

	switch v := m.Value.(type) {
	case bthomev2_types.NumberValue:
		rec.Value = new(v.Number)
	case bthomev2_types.BinaryValue:
		rec.BoolValue = new(v.Boolean)
	case bthomev2_types.TextValue:
		val := v.Text
		rec.StringValue = &val
	case bthomev2_types.RawValue:
		rec.DataValue = new(string(v.Raw))
	case bthomev2_types.EventValue:
		rec.StringValue = new(eventToSenML[v.Event].ToString())
	}

	return rec
}

func UnitToSenML(p bthomev2_types.Property, u bthomev2_types.Unit) senmltypes.HubUnit {
	switch u {
	case bthomev2_types.CelsiusDegree:
		return senmltypes.CelsiusDegree
	case bthomev2_types.Percentage:
		if p == bthomev2_types.SensorHumidity {
			return senmltypes.HumidityPercentage
		}
		return senmltypes.Percentage
	case bthomev2_types.Volt:
		return senmltypes.Volt
	case bthomev2_types.Ampere:
		return senmltypes.Ampere
	case bthomev2_types.Watt:
		return senmltypes.Watt
	case bthomev2_types.KilowattHours:
		return senmltypes.KilowattHours
	case bthomev2_types.MetersPerSecondSquared:
		return senmltypes.MetersPerSecondSquared
	case bthomev2_types.PartsPerMillion:
		return senmltypes.PartsPerMillion
	case bthomev2_types.Microsiemens:
		return senmltypes.Microsiemens
	case bthomev2_types.Degree:
		return senmltypes.Degree
	case bthomev2_types.Millimeter:
		return senmltypes.Millimeter
	case bthomev2_types.Meter:
		return senmltypes.Meter
	case bthomev2_types.Second:
		return senmltypes.Second
	case bthomev2_types.CubicMeter:
		return senmltypes.CubicMeter
	case bthomev2_types.DegreesPerSecond:
		return senmltypes.DegreesPerSecond
	case bthomev2_types.Lux:
		return senmltypes.Lux
	case bthomev2_types.Kilogram:
		return senmltypes.Kilogram
	case bthomev2_types.Pound:
		return senmltypes.Pound
	case bthomev2_types.MicrogramsPerCubicMeter:
		return senmltypes.MicrogramsPerCubicMeter
	case bthomev2_types.Hectopascal:
		return senmltypes.Hectopascal
	case bthomev2_types.RevolutionsPerMinute:
		return senmltypes.RevolutionsPerMinute
	case bthomev2_types.MetersPerSecond:
		return senmltypes.MetersPerSecond
	case bthomev2_types.Liter:
		return senmltypes.Liter
	case bthomev2_types.Milliliter:
		return senmltypes.Milliliter
	case bthomev2_types.CubicMetersPerHour:
		return senmltypes.CubicMetersPerHour
	case bthomev2_types.Boolean:
		return senmltypes.Empty
	default:
		return senmltypes.Empty
	}
}

var propertyToName = map[bthomev2_types.Property]senmltypes.HubRecordName{
	// --- Numeric values ---
	bthomev2_types.SensorAcceleration:  senmltypes.HubRecordAcceleration,
	bthomev2_types.SensorBattery:       senmltypes.HubRecordBattery,
	bthomev2_types.SensorChannel:       senmltypes.HubRecordChannel,
	bthomev2_types.SensorCO2:           senmltypes.HubRecordCO2,
	bthomev2_types.SensorConductivity:  senmltypes.HubRecordConductivity,
	bthomev2_types.SensorCount:         senmltypes.HubRecordCount,
	bthomev2_types.SensorCurrent:       senmltypes.HubRecordCurrent,
	bthomev2_types.SensorDewPoint:      senmltypes.HubRecordDewPoint,
	bthomev2_types.SensorDirection:     senmltypes.HubRecordDirection,
	bthomev2_types.SensorDistanceMM:    senmltypes.HubRecordDistance,
	bthomev2_types.SensorDistanceM:     senmltypes.HubRecordDistance,
	bthomev2_types.SensorDuration:      senmltypes.HubRecordDuration,
	bthomev2_types.SensorEnergy:        senmltypes.HubRecordEnergy,
	bthomev2_types.SensorGas:           senmltypes.HubRecordGas,
	bthomev2_types.SensorGyroscope:     senmltypes.HubRecordGyroscope,
	bthomev2_types.SensorHumidity:      senmltypes.HubRecordHumidity,
	bthomev2_types.SensorIlluminance:   senmltypes.HubRecordIlluminance,
	bthomev2_types.SensorMassKG:        senmltypes.HubRecordMass,
	bthomev2_types.SensorMassLB:        senmltypes.HubRecordMass,
	bthomev2_types.SensorMoisture:      senmltypes.HubRecordMoisture,
	bthomev2_types.SensorPM2_5:         senmltypes.HubRecordPM2_5,
	bthomev2_types.SensorPM10:          senmltypes.HubRecordPM10,
	bthomev2_types.SensorPower:         senmltypes.HubRecordPower,
	bthomev2_types.SensorPrecipitation: senmltypes.HubRecordPrecipitation,
	bthomev2_types.SensorPressure:      senmltypes.HubRecordPressure,
	bthomev2_types.SensorRaw:           senmltypes.HubRecordRaw,
	bthomev2_types.SensorRotation:      senmltypes.HubRecordRotation,
	bthomev2_types.SensorRotational:    senmltypes.HubRecordRotational,
	bthomev2_types.SensorSpeed:         senmltypes.HubRecordSpeed,
	bthomev2_types.SensorTemperature:   senmltypes.HubRecordTemperature,
	bthomev2_types.SensorText:          senmltypes.HubRecordText,
	bthomev2_types.SensorTimestamp:     senmltypes.HubRecordTimestamp,
	bthomev2_types.SensorTVOC:          senmltypes.HubRecordTVOC,
	bthomev2_types.SensorVoltage:       senmltypes.HubRecordVoltage,
	bthomev2_types.SensorVolume:        senmltypes.HubRecordVolume,
	bthomev2_types.SensorVolumeML:      senmltypes.HubRecordVolume,
	bthomev2_types.SensorVolumeStorage: senmltypes.HubRecordVolume,
	bthomev2_types.SensorVolumeFlow:    senmltypes.HubRecordVolumeFlow,
	bthomev2_types.SensorUV:            senmltypes.HubRecordUV,
	bthomev2_types.SensorWater:         senmltypes.HubRecordWater,

	// --- Bool values / States ---
	bthomev2_types.SensorBatteryCharging: senmltypes.HubRecordBatteryCharging,
	bthomev2_types.SensorCarbonMonoxide:  senmltypes.HubRecordCarbonMonoxide,
	bthomev2_types.SensorCold:            senmltypes.HubRecordCold,
	bthomev2_types.SensorConnectivity:    senmltypes.HubRecordConnectivity,
	bthomev2_types.SensorDoor:            senmltypes.HubRecordDoor,
	bthomev2_types.SensorGarageDoor:      senmltypes.HubRecordGarageDoor,
	bthomev2_types.SensorGenericBoolean:  senmltypes.HubRecordGenericBoolean,
	bthomev2_types.SensorHeat:            senmltypes.HubRecordHeat,
	bthomev2_types.SensorLight:           senmltypes.HubRecordLight,
	bthomev2_types.SensorLock:            senmltypes.HubRecordLock,
	bthomev2_types.SensorMotion:          senmltypes.HubRecordMotion,
	bthomev2_types.SensorMoving:          senmltypes.HubRecordMoving,
	bthomev2_types.SensorOccupancy:       senmltypes.HubRecordOccupancy,
	bthomev2_types.SensorOpening:         senmltypes.HubRecordOpening,
	bthomev2_types.SensorPlug:            senmltypes.HubRecordPlug,
	bthomev2_types.SensorPresence:        senmltypes.HubRecordPresence,
	bthomev2_types.SensorProblem:         senmltypes.HubRecordProblem,
	bthomev2_types.SensorRunning:         senmltypes.HubRecordRunning,
	bthomev2_types.SensorSafety:          senmltypes.HubRecordSafety,
	bthomev2_types.SensorSmoke:           senmltypes.HubRecordSmoke,
	bthomev2_types.SensorSound:           senmltypes.HubRecordSound,
	bthomev2_types.SensorTamper:          senmltypes.HubRecordTamper,
	bthomev2_types.SensorVibration:       senmltypes.HubRecordVibration,
	bthomev2_types.SensorWindow:          senmltypes.HubRecordWindow,

	// --- Events ---
	bthomev2_types.SensorButtonEvent: senmltypes.HubRecordButtonEvent,
	bthomev2_types.SensorDimmerEvent: senmltypes.HubRecordDimmerEvent,
}

var eventToSenML = map[bthomev2_types.Event]senmltypes.HubEvent{
	bthomev2_types.EventButtonNone:            senmltypes.HubEventNone,
	bthomev2_types.EventButtonPress:           senmltypes.HubEventButtonPress,
	bthomev2_types.EventButtonDoublePress:     senmltypes.HubEventButtonDoublePress,
	bthomev2_types.EventButtonTriplePress:     senmltypes.HubEventButtonTriplePress,
	bthomev2_types.EventButtonLongPress:       senmltypes.HubEventButtonLongPress,
	bthomev2_types.EventButtonLongDoublePress: senmltypes.HubEventButtonLongDoublePress,
	bthomev2_types.EventButtonLongTriplePress: senmltypes.HubEventButtonLongTriplePress,
	bthomev2_types.EventButtonHoldPress:       senmltypes.HubEventButtonHoldPress,
	bthomev2_types.EventDimmerRotateLeft:      senmltypes.HubEventDimmerRotateLeft,
	bthomev2_types.EventDimmerRotateRight:     senmltypes.HubEventDimmerRotateRight,
}
