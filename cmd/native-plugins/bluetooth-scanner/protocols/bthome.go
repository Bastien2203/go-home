package protocols

import (
	"fmt"

	"log"

	"github.com/Bastien2203/bthomev2"
	"github.com/Bastien2203/go-home/shared/types"
	"tinygo.org/x/bluetooth"
)

var BthomeUUID = bluetooth.New16BitUUID(uint16(bthomev2.ServiceDataUUID))

type BthomeParser struct{}

func NewBthomeParser() *BthomeParser {
	return &BthomeParser{}
}

func (d *BthomeParser) Name() string {
	return "bthome"
}

func (d *BthomeParser) CanParse() bool {
	return true
}

func (d *BthomeParser) Parse(payload []byte) ([]*types.Capability, error) {
	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, fmt.Errorf("error while parsing service data: %w", err)
	}

	capabilities := make([]*types.Capability, 0)
	for _, m := range data {
		switch v := m.Value.(type) {
		case bthomev2.NumberValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Number,
				Type:  types.TypeFloat,
				Unit:  types.UnitFromBthome(m.Unit),
			})
		case bthomev2.BinaryValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Boolean,
				Type:  types.TypeBool,
				Unit:  types.UnitFromBthome(m.Unit),
			})
		case bthomev2.TextValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Text,
				Type:  types.TypeString,
				Unit:  types.UnitFromBthome(m.Unit),
			})
		case bthomev2.RawValue:
			capabilities = append(capabilities, &types.Capability{
				Name:  types.CapabilityType(m.Property),
				Value: v.Raw,
				Type:  types.TypeBytes,
				Unit:  types.UnitFromBthome(m.Unit),
			})
		}
	}

	return capabilities, nil
}
