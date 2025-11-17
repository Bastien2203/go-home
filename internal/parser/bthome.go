package parser

import (
	"gohome/internal/core"
	"gohome/internal/scanner"
	"log"

	"github.com/Bastien2203/bthomev2"
	"tinygo.org/x/bluetooth"
)

type BthomeParser struct {
	scanner *scanner.BluetoothScanner
}

var serviceDataUUID = bluetooth.New16BitUUID(uint16(bthomev2.ServiceDataUUID))

func NewBthomeParser(scanner *scanner.BluetoothScanner) core.Parser {
	return &BthomeParser{scanner: scanner}
}

func (p *BthomeParser) Scanner() core.Scanner { return p.scanner }

func (p *BthomeParser) Name() string { return "bthome" }

func (p *BthomeParser) Parse(adv core.Advertisment) (map[string]any, bool) {
	bluetoothAdv, ok := adv.(*core.BluetoothAdvertisement)
	if !ok {
		return nil, false
	}
	payload, ok := bluetoothAdv.ServiceData[serviceDataUUID]
	if !ok {
		return nil, false
	}

	data, err := bthomev2.ParseServiceData(payload)
	if err != nil {
		log.Printf("error while parsing service data : %s\n", err.Error())
		return nil, false
	}

	measurements := make(map[string]any, len(data))
	for _, m := range data {
		switch v := m.Value.(type) {
		case bthomev2.NumberValue:
			measurements[string(m.Property)] = v.Number
		case bthomev2.BinaryValue:
			measurements[string(m.Property)] = v.Boolean
		case bthomev2.TextValue:
			measurements[string(m.Property)] = v.Text
		case bthomev2.RawValue:
			measurements[string(m.Property)] = v.Raw
		default:
			measurements[string(m.Property)] = v // fallback, should not happen
		}
	}

	return measurements, true
}
