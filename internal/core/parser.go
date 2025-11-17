package core

import "time"

type Parser interface {
	Name() string
	Parse(adv Advertisment) (map[string]any, bool)
	Scanner() Scanner
}

type ParsedMessage struct {
	Addr     string
	DeviceID string
	Values   map[string]any
	Time     time.Time
}
