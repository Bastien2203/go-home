package types

type Protocol interface {
	ID() string
	Name() string
	Parse(raw []byte) ([]*Capability, error)
	AddressType() AddressType
}
