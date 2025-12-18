package main

import (
	"bluetooth-scanner/protocols"

	"github.com/Bastien2203/go-home/shared/types"
)

type Protocol interface {
	Name() string
	Parse(payload []byte) ([]*types.Capability, error)
}

var ProtocolList = map[string]Protocol{
	protocols.BthomeUUID: protocols.NewBthomeParser(),
}
