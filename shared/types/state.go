package types

type State string

const (
	StateRunning    State = "running"
	StateRestarting State = "restarting"
	StateStopped    State = "stopped"
)
