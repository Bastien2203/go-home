package protocols

import (
	"testing"
)

func TestSwitchBotParser(t *testing.T) {
	address := "364f1e8a-0a7d-56ca-d70a-e549b246286f"
	payload := []byte{0x54, 0x00, 0xe4, 0x07, 0x92, 0x3a}

	parser := NewSwitchBotParser()

	_, _, err := parser.Parse(address, payload)
	if err != nil {
		t.Fatalf("Failed to parse payload %v", err)
	}
}
