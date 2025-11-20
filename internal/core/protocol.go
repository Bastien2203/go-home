package core

type Protocol interface {
	ID() string
	Name() string
	Parse(raw any) ([]*Capability, error)
	AddressType() AddressType
}
