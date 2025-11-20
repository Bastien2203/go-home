package protocols

import (
	"fmt"
	"gohome/internal/core"
	"log"

	"github.com/Bastien2203/bthomev2"
	"tinygo.org/x/bluetooth"
)

var serviceDataUUID = bluetooth.New16BitUUID(uint16(bthomev2.ServiceDataUUID))

type BthomeParser struct{}

func NewBthomeParser() *BthomeParser {
	return &BthomeParser{}
}

func (d *BthomeParser) AddressType() core.AddressType {
	return core.BLEAddress
}

func (s *BthomeParser) ID() string {
	return "bthome"
}

func (d *BthomeParser) Name() string {
	return "BTHome"
}

func (d *BthomeParser) Parse(rawData any) ([]*core.Capability, error) {
	bluetoothAdv, ok := rawData.(core.BluetoothAdvertisement)
	if !ok {
		return nil, fmt.Errorf("invalid advertisement type")
	}

	payload, ok := bluetoothAdv[serviceDataUUID]
	if !ok {
		return nil, fmt.Errorf("service data not found")
	}

	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, fmt.Errorf("error while parsing service data: %w", err)
	}

	capabilities := make([]*core.Capability, 0)
	for _, m := range data {
		switch v := m.Value.(type) {
		case bthomev2.NumberValue:
			capabilities = append(capabilities, &core.Capability{
				Name:  core.CapabilityType(m.Property),
				Value: v.Number,
				Type:  core.TypeFloat,
				Unit:  core.UnitFromBthome(m.Unit),
			})
		case bthomev2.BinaryValue:
			capabilities = append(capabilities, &core.Capability{
				Name:  core.CapabilityType(m.Property),
				Value: v.Boolean,
				Type:  core.TypeBool,
				Unit:  core.UnitFromBthome(m.Unit),
			})
		case bthomev2.TextValue:
			capabilities = append(capabilities, &core.Capability{
				Name:  core.CapabilityType(m.Property),
				Value: v.Text,
				Type:  core.TypeString,
				Unit:  core.UnitFromBthome(m.Unit),
			})
		case bthomev2.RawValue:
			capabilities = append(capabilities, &core.Capability{
				Name:  core.CapabilityType(m.Property),
				Value: v.Raw,
				Type:  core.TypeBytes,
				Unit:  core.UnitFromBthome(m.Unit),
			})
		}
	}

	return capabilities, nil
}
