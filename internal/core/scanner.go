package core

import "context"

type Scanner interface {
	ID() string
	Name() string
	Start(ctx context.Context) error
	Stop() error
	State() State
}
