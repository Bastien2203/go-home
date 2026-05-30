package protocols

import (
	"fmt"

	"github.com/absmach/senml"
)

type NotImplementedParser struct {
	name string
}

func NewNotImplementedParser(name string) *NotImplementedParser {
	return &NotImplementedParser{name: name}
}

func (d *NotImplementedParser) Name() string {
	return d.name
}

func (d *NotImplementedParser) CanParse() bool {
	return false
}

func (d *NotImplementedParser) Parse(address string, payload []byte) ([]senml.Record, bool, error) {
	return nil, false, fmt.Errorf("not implemented")
}
