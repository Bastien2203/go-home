package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

const HostMACAddress = "a8:51:ab:bb:45:5b"

func toggleSpeaker(adapter *bluetooth.Adapter, deviceAddress bluetooth.Address, turnOn bool) error {
	fmt.Println("Connecting to speaker...")

	device, err := adapter.Connect(deviceAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return fmt.Errorf("failed to connect: %v", err)
	}
	defer device.Disconnect()
	fmt.Println("Connected!")

	ueServiceUUID, _ := bluetooth.ParseUUID("000061fe-0000-1000-8000-00805f9b34fb")
	services, err := device.DiscoverServices([]bluetooth.UUID{ueServiceUUID})
	if err != nil || len(services) == 0 {
		return fmt.Errorf("failed to find UE service: %v", err)
	}

	powerCharUUID, _ := bluetooth.ParseUUID("c6d6dc0d-07f5-47ef-9b59-630622b01fd3")
	chars, err := services[0].DiscoverCharacteristics([]bluetooth.UUID{powerCharUUID})
	if err != nil || len(chars) == 0 {
		return fmt.Errorf("failed to find power characteristic: %v", err)
	}

	payload, err := buildPayload(HostMACAddress, turnOn)
	if err != nil {
		return fmt.Errorf("failed to build payload: %v", err)
	}

	fmt.Printf("Sending payload: %X\n", payload)
	_, err = chars[0].WriteWithoutResponse(payload)
	if err != nil {
		return fmt.Errorf("failed to write command: %v", err)
	}

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Command sent successfully!")
	return nil
}

func buildPayload(macStr string, turnOn bool) ([]byte, error) {
	cleanMac := strings.ReplaceAll(macStr, ":", "")

	macBytes, err := hex.DecodeString(cleanMac)
	if err != nil || len(macBytes) != 6 {
		return nil, fmt.Errorf("invalid MAC address format")
	}

	cmdByte := byte(0x01) // ON
	if !turnOn {
		cmdByte = byte(0x02) // OFF
	}

	// Append command byte to the end of the MAC bytes
	return append(macBytes, cmdByte), nil
}
