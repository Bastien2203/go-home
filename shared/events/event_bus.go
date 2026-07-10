package events

import (
	"encoding/json"
	"fmt"
	"strings"

	"log"
	"time"

	"github.com/Bastien2203/go-home/shared/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fxamacker/cbor/v2"
)

type EventType string

const (
	ParsedDataReceived    EventType = "gohome/parsed_data"
	BluetoothDeviceFound  EventType = "gohome/bluetooth/found"
	PluginConnected       EventType = "gohome/plugin/connected"
	PluginDisconnected    EventType = "gohome/plugin/disconnected"
	PluginStateChanged    EventType = "gohome/plugin/newstate"
	CoreDiscovery         EventType = "gohome/core/discovery"
	PluginCommandResponse EventType = "gohome/plugin/command-response"
)

func PluginStop(id string) EventType {
	return EventType(fmt.Sprintf("gohome/plugin/stop/%s", id))
}

func PluginStart(id string) EventType {
	return EventType(fmt.Sprintf("gohome/plugin/start/%s", id))
}

func RegisterDeviceForAdapter(id string) EventType {
	return EventType(fmt.Sprintf("gohome/device/register/%s", id))
}

func UnregisterDeviceForAdapter(id string) EventType {
	return EventType(fmt.Sprintf("gohome/device/unregister/%s", id))
}

func UpdateDataForAdapter(id string) EventType {
	return EventType(fmt.Sprintf("gohome/device/updated/%s", id))
}

type Event struct {
	Type    EventType
	Payload any
}

type EventBus struct {
	client mqtt.Client
	debug  bool
}

func NewEventBus(brokerURL string, clientID string) (*EventBus, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURL)
	opts.SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Println("[EventBus] connected to the mqtt broker")
	})
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Println("[EventBus] connection to mqtt broker lost")
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return &EventBus{client: client, debug: config.IsDebug()}, nil
}

func (eb *EventBus) subscribeRaw(eventType EventType, handler func(payload []byte)) error {
	topic := string(eventType)

	mqttHandler := func(client mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	}

	token := eb.client.Subscribe(topic, 0, mqttHandler)
	token.Wait()
	return token.Error()
}

var cborTopicPrefixes = []string{
	string(ParsedDataReceived),
	"gohome/device/updated",
}

func isCBORTopic(eventType EventType) bool {
	topic := string(eventType)
	for _, prefix := range cborTopicPrefixes {
		if topic == prefix || strings.HasPrefix(topic, prefix+"/") {
			return true
		}
	}
	return false
}

func Subscribe[T any](eb *EventBus, eventType EventType, handler func(T)) error {
	return eb.subscribeRaw(eventType, func(rawPayload []byte) {
		var target T
		var err error

		if isCBORTopic(eventType) {
			err = cbor.Unmarshal(rawPayload, &target)
		} else {
			err = json.Unmarshal(rawPayload, &target)
		}

		if err != nil {
			log.Printf("Unmarshal error on topic %s: %v", eventType, err)
			return
		}

		handler(target)
	})
}

func (eb *EventBus) Publish(event Event) error {
	var payloadBytes []byte
	var err error

	if isCBORTopic(event.Type) {
		payloadBytes, err = cbor.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload to CBOR: %w", err)
		}

		if eb.debug {
			log.Printf("[EventBus] Publish event (%s) : %d bytes encoded in CBOR\n", event.Type, len(payloadBytes))
		}
	} else {
		payloadBytes, err = json.Marshal(event.Payload)
		if err != nil {
			return fmt.Errorf("failed to marshal payload: %w", err)
		}

		if eb.debug {
			log.Printf("[EventBus] Publish event (%s) : %s\n", event.Type, string(payloadBytes))
		}
	}

	token := eb.client.Publish(string(event.Type), 0, false, payloadBytes)
	token.Wait()

	return token.Error()
}

func (eb *EventBus) Close() {
	eb.client.Disconnect(250)
}
