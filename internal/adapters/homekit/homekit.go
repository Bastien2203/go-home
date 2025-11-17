package homekit

import (
	"context"
	"errors"
	"fmt"
	"gohome/internal/core"
	"log"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

type hkMeta struct {
	Acc        *accessory.A
	DeviceAddr string
}
type HomeKitAdapter struct {
	*core.BaseAdapter
	server *hap.Server
}

func NewHomeKitAdapter() core.Adapter {
	h := &HomeKitAdapter{}
	base := core.NewBaseAdapter("HomeKit Adapter", h)
	h.BaseAdapter = base
	return h
}

func (h *HomeKitAdapter) OnStart(ctx context.Context) error {
	// collect accessories already registered
	metas := h.DeviceMetas()
	if len(metas) == 0 {
		return errors.New("homekit adapter needs at least one accessory to start")
	}

	fs := hap.NewFsStore("./db")

	hapAccs := make([]*accessory.A, 0, len(metas))
	for _, m := range metas {
		if hk, ok := m.(hkMeta); ok {
			hapAccs = append(hapAccs, hk.Acc)
		}
	}

	server, err := hap.NewServer(fs, hapAccs[0], hapAccs[1:]...)
	if err != nil {
		return err
	}
	h.server = server

	go func() {
		if err := server.ListenAndServe(ctx); err != nil {
			log.Printf("homekit server error: %v", err)
		}
	}()
	return nil
}

func (h *HomeKitAdapter) OnStop() error {
	// hap.Server does not expose Stop; we cancelled base ctx in BaseAdapter.Stop()
	return nil
}

func (h *HomeKitAdapter) OnRegisterDevice(d core.Device) (any, context.CancelFunc, error) {
	ctx := h.Context()
	var a *accessory.A
	var cancel context.CancelFunc

	switch d.Type() {
	case core.TemperatureSensorType:
		acc := accessory.NewTemperatureSensor(accessory.Info{Name: d.Name()})
		a = acc.A
		deviceCtx, c := context.WithCancel(ctx)
		cancel = c

		go func() {
			events := make(chan core.Event, 8)
			d.Subscribe(events)
			defer d.Unsubscribe(events)

			for {
				select {
				case ev := <-events:
					if ev.Type == "state_change" {
						if v, ok := ev.Payload["temperature"].(float64); ok {
							acc.TempSensor.CurrentTemperature.SetValue(v)
						}
					}
				case <-deviceCtx.Done():
					return
				}
			}
		}()
	default:
		return nil, nil, fmt.Errorf("unsupported device type: %s", d.Type())
	}

	meta := hkMeta{Acc: a, DeviceAddr: d.Addr()}
	return meta, cancel, nil
}

func (h *HomeKitAdapter) OnUnregisterDevice(meta any) error {
	// adapter-specific cleanup if needed
	// base already called cancel()
	// nothing else to do for now
	return nil
}
