package protocols

import (
	"fmt"
	"time"

	"log"

	"github.com/Bastien2203/bthomev2"
	bthomev2_types "github.com/Bastien2203/bthomev2/types"
	"github.com/Bastien2203/go-home/shared/types"

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

func UnitToSenML(p bthomev2_types.Property, u bthomev2_types.Unit) types.HubUnit {
	switch u {
	case bthomev2_types.CelsiusDegree:
		return types.CelsiusDegree
	case bthomev2_types.Percentage:
		if p == bthomev2_types.SensorHumidity {
			return types.HumidityPercentage
		}
		return types.Percentage
	case bthomev2_types.Volt:
		return types.Volt
	case bthomev2_types.Ampere:
		return types.Ampere
	case bthomev2_types.Watt:
		return types.Watt
	case bthomev2_types.KilowattHours:
		return types.KilowattHours
	case bthomev2_types.MetersPerSecondSquared:
		return types.MetersPerSecondSquared
	case bthomev2_types.PartsPerMillion:
		return types.PartsPerMillion
	case bthomev2_types.Microsiemens:
		return types.Microsiemens
	case bthomev2_types.Degree:
		return types.Degree
	case bthomev2_types.Millimeter:
		return types.Millimeter
	case bthomev2_types.Meter:
		return types.Meter
	case bthomev2_types.Second:
		return types.Second
	case bthomev2_types.CubicMeter:
		return types.CubicMeter
	case bthomev2_types.DegreesPerSecond:
		return types.DegreesPerSecond
	case bthomev2_types.Lux:
		return types.Lux
	case bthomev2_types.Kilogram:
		return types.Kilogram
	case bthomev2_types.Pound:
		return types.Pound
	case bthomev2_types.MicrogramsPerCubicMeter:
		return types.MicrogramsPerCubicMeter
	case bthomev2_types.Hectopascal:
		return types.Hectopascal
	case bthomev2_types.RevolutionsPerMinute:
		return types.RevolutionsPerMinute
	case bthomev2_types.MetersPerSecond:
		return types.MetersPerSecond
	case bthomev2_types.Liter:
		return types.Liter
	case bthomev2_types.Milliliter:
		return types.Milliliter
	case bthomev2_types.CubicMetersPerHour:
		return types.CubicMetersPerHour
	case bthomev2_types.Boolean:
		return types.Empty
	default:
		return types.Empty
	}
}

var propertyToName = map[bthomev2_types.Property]types.HubRecordName{
	// --- Numeric values ---
	bthomev2_types.SensorAcceleration:  types.HubRecordAcceleration,
	bthomev2_types.SensorBattery:       types.HubRecordBattery,
	bthomev2_types.SensorChannel:       types.HubRecordChannel,
	bthomev2_types.SensorCO2:           types.HubRecordCO2,
	bthomev2_types.SensorConductivity:  types.HubRecordConductivity,
	bthomev2_types.SensorCount:         types.HubRecordCount,
	bthomev2_types.SensorCurrent:       types.HubRecordCurrent,
	bthomev2_types.SensorDewPoint:      types.HubRecordDewPoint,
	bthomev2_types.SensorDirection:     types.HubRecordDirection,
	bthomev2_types.SensorDistanceMM:    types.HubRecordDistance,
	bthomev2_types.SensorDistanceM:     types.HubRecordDistance,
	bthomev2_types.SensorDuration:      types.HubRecordDuration,
	bthomev2_types.SensorEnergy:        types.HubRecordEnergy,
	bthomev2_types.SensorGas:           types.HubRecordGas,
	bthomev2_types.SensorGyroscope:     types.HubRecordGyroscope,
	bthomev2_types.SensorHumidity:      types.HubRecordHumidity,
	bthomev2_types.SensorIlluminance:   types.HubRecordIlluminance,
	bthomev2_types.SensorMassKG:        types.HubRecordMass,
	bthomev2_types.SensorMassLB:        types.HubRecordMass,
	bthomev2_types.SensorMoisture:      types.HubRecordMoisture,
	bthomev2_types.SensorPM2_5:         types.HubRecordPM2_5,
	bthomev2_types.SensorPM10:          types.HubRecordPM10,
	bthomev2_types.SensorPower:         types.HubRecordPower,
	bthomev2_types.SensorPrecipitation: types.HubRecordPrecipitation,
	bthomev2_types.SensorPressure:      types.HubRecordPressure,
	bthomev2_types.SensorRaw:           types.HubRecordRaw,
	bthomev2_types.SensorRotation:      types.HubRecordRotation,
	bthomev2_types.SensorRotational:    types.HubRecordRotational,
	bthomev2_types.SensorSpeed:         types.HubRecordSpeed,
	bthomev2_types.SensorTemperature:   types.HubRecordTemperature,
	bthomev2_types.SensorText:          types.HubRecordText,
	bthomev2_types.SensorTimestamp:     types.HubRecordTimestamp,
	bthomev2_types.SensorTVOC:          types.HubRecordTVOC,
	bthomev2_types.SensorVoltage:       types.HubRecordVoltage,
	bthomev2_types.SensorVolume:        types.HubRecordVolume,
	bthomev2_types.SensorVolumeML:      types.HubRecordVolume,
	bthomev2_types.SensorVolumeStorage: types.HubRecordVolume,
	bthomev2_types.SensorVolumeFlow:    types.HubRecordVolumeFlow,
	bthomev2_types.SensorUV:            types.HubRecordUV,
	bthomev2_types.SensorWater:         types.HubRecordWater,

	// --- Bool values / States ---
	bthomev2_types.SensorBatteryCharging: types.HubRecordBatteryCharging,
	bthomev2_types.SensorCarbonMonoxide:  types.HubRecordCarbonMonoxide,
	bthomev2_types.SensorCold:            types.HubRecordCold,
	bthomev2_types.SensorConnectivity:    types.HubRecordConnectivity,
	bthomev2_types.SensorDoor:            types.HubRecordDoor,
	bthomev2_types.SensorGarageDoor:      types.HubRecordGarageDoor,
	bthomev2_types.SensorGenericBoolean:  types.HubRecordGenericBoolean,
	bthomev2_types.SensorHeat:            types.HubRecordHeat,
	bthomev2_types.SensorLight:           types.HubRecordLight,
	bthomev2_types.SensorLock:            types.HubRecordLock,
	bthomev2_types.SensorMotion:          types.HubRecordMotion,
	bthomev2_types.SensorMoving:          types.HubRecordMoving,
	bthomev2_types.SensorOccupancy:       types.HubRecordOccupancy,
	bthomev2_types.SensorOpening:         types.HubRecordOpening,
	bthomev2_types.SensorPlug:            types.HubRecordPlug,
	bthomev2_types.SensorPresence:        types.HubRecordPresence,
	bthomev2_types.SensorProblem:         types.HubRecordProblem,
	bthomev2_types.SensorRunning:         types.HubRecordRunning,
	bthomev2_types.SensorSafety:          types.HubRecordSafety,
	bthomev2_types.SensorSmoke:           types.HubRecordSmoke,
	bthomev2_types.SensorSound:           types.HubRecordSound,
	bthomev2_types.SensorTamper:          types.HubRecordTamper,
	bthomev2_types.SensorVibration:       types.HubRecordVibration,
	bthomev2_types.SensorWindow:          types.HubRecordWindow,

	// --- Events ---
	bthomev2_types.SensorButtonEvent: types.HubRecordButtonEvent,
	bthomev2_types.SensorDimmerEvent: types.HubRecordDimmerEvent,
}

var eventToSenML = map[bthomev2_types.Event]types.HubEvent{
	bthomev2_types.EventButtonNone:            types.HubEventNone,
	bthomev2_types.EventButtonPress:           types.HubEventButtonPress,
	bthomev2_types.EventButtonDoublePress:     types.HubEventButtonDoublePress,
	bthomev2_types.EventButtonTriplePress:     types.HubEventButtonTriplePress,
	bthomev2_types.EventButtonLongPress:       types.HubEventButtonLongPress,
	bthomev2_types.EventButtonLongDoublePress: types.HubEventButtonLongDoublePress,
	bthomev2_types.EventButtonLongTriplePress: types.HubEventButtonLongTriplePress,
	bthomev2_types.EventButtonHoldPress:       types.HubEventButtonHoldPress,
	bthomev2_types.EventDimmerRotateLeft:      types.HubEventDimmerRotateLeft,
	bthomev2_types.EventDimmerRotateRight:     types.HubEventDimmerRotateRight,
}
