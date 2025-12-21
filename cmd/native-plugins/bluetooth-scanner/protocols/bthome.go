package protocols

import (
	"bluetooth-scanner/utils"
	"fmt"

	"log"

	"github.com/Bastien2203/bthomev2"
	bthomev2_types "github.com/Bastien2203/bthomev2/types"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

var BthomeUUID = bluetooth.New16BitUUID(uint16(bthomev2_types.ServiceDataUUID))

type BthomeParser struct {
	cache map[string]uint8
}

func NewBthomeParser() *BthomeParser {
	return &BthomeParser{
		cache: make(map[string]uint8),
	}
}

func (d *BthomeParser) Name() string {
	return "bthome"
}

func (d *BthomeParser) CanParse() bool {
	return true
}

func (d *BthomeParser) Parse(address string, payload []byte) ([]*types.Capability, error) {
	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, fmt.Errorf("error while parsing service data: %w", err)
	}

	pid, found := data[bthomev2_types.PacketID]
	if found {
		pidValue, ok := pid.Float64()
		if !ok {
			return nil, fmt.Errorf("invalid packet id value")
		}
		// If packet already seen, ignore it
		if entry, ok := d.cache[address]; ok {
			if entry == uint8(pidValue) {
				return nil, utils.DeduplicateBluetoothPackets
			}
		}
		// Else, update cache
		d.cache[address] = uint8(pidValue)
	}

	capabilities := make([]*types.Capability, 0)
	for _, m := range data {
		switch v := m.Value.(type) {
		case bthomev2_types.NumberValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Number,
				Type:  types.TypeFloat,
				Unit:  UnitFromBthome(m.Unit),
			})
		case bthomev2_types.BinaryValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Boolean,
				Type:  types.TypeBool,
				Unit:  UnitFromBthome(m.Unit),
			})
		case bthomev2_types.TextValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Text,
				Type:  types.TypeString,
				Unit:  UnitFromBthome(m.Unit),
			})
		case bthomev2_types.RawValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Raw,
				Type:  types.TypeBytes,
				Unit:  UnitFromBthome(m.Unit),
			})

		case bthomev2_types.EventValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Event,
				Unit:  UnitFromBthome(m.Unit),
			})
		}
	}

	return capabilities, nil
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
