package protocols

import (
	"fmt"

	"github.com/Bastien2203/go-home/shared/types"
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

func (d *NotImplementedParser) Parse(address string, payload []byte) ([]*types.Capability, bool, error) {
	return nil, false, fmt.Errorf("not implemented")
}
