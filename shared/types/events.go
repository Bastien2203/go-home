package types

type HubEvent string

func (e HubEvent) ToString() string {
	return string(e)
}

const (
	HubEventNone                  HubEvent = "none"
	HubEventButtonPress           HubEvent = "press"
	HubEventButtonDoublePress     HubEvent = "double_press"
	HubEventButtonTriplePress     HubEvent = "triple_press"
	HubEventButtonLongPress       HubEvent = "long_press"
	HubEventButtonLongDoublePress HubEvent = "long_double_press"
	HubEventButtonLongTriplePress HubEvent = "long_triple_press"
	HubEventButtonHoldPress       HubEvent = "hold_press"
	HubEventDimmerRotateLeft      HubEvent = "rotate_left"
	HubEventDimmerRotateRight     HubEvent = "rotate_right"
	HubEventUnknown               HubEvent = "unknown"
)
