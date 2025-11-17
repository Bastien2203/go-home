package core

import (
	"context"
	"errors"
	"gohome/utils"
	"sync"

	"github.com/google/uuid"
)

type Adapter interface {
	Id() string
	Name() string
	Start(ctx context.Context) error
	Stop() error
	RegisterDevice(d Device) error
	UnregisterDevice(addr string) error
	State() string
	ToJson() map[string]any
}

type RegisteredDevice struct {
	Device Device
	Meta   any
	Cancel context.CancelFunc
}

type AdapterHooks interface {
	OnStart(ctx context.Context) error
	OnStop() error
	OnRegisterDevice(d Device) (meta any, cancel context.CancelFunc, err error)
	OnUnregisterDevice(meta any) error
}

type BaseAdapter struct {
	name  string
	hooks AdapterHooks

	mu      sync.RWMutex
	devices map[string]RegisteredDevice

	ctx         context.Context
	cancel      context.CancelFunc
	needRestart bool
	id          string
}

func NewBaseAdapter(name string, hooks AdapterHooks) *BaseAdapter {
	return &BaseAdapter{
		id:          uuid.New().String(),
		name:        name,
		hooks:       hooks,
		devices:     make(map[string]RegisteredDevice),
		needRestart: false,
	}
}

func (b *BaseAdapter) Id() string { return b.id }

func (b *BaseAdapter) Name() string { return b.name }

func (b *BaseAdapter) Start(ctx context.Context) error {
	b.mu.Lock()
	b.needRestart = false
	if b.cancel != nil {
		b.mu.Unlock()
		return errors.New("adapter already started")
	}
	b.ctx, b.cancel = context.WithCancel(ctx)
	b.mu.Unlock()

	// restart linked devices :
	for _, rd := range b.devices {
		d := rd.Device
		meta, cancel, err := b.hooks.OnRegisterDevice(d)
		if err != nil {
			return err
		}

		// store registration
		b.mu.Lock()
		b.devices[d.Addr()] = RegisteredDevice{
			Device: d,
			Meta:   meta,
			Cancel: cancel,
		}
		b.mu.Unlock()
	}

	return b.hooks.OnStart(b.Context())
}

func (b *BaseAdapter) State() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.cancel != nil {
		if b.needRestart {
			return "need_restart"
		}
		return "started"
	}
	return "stopped"
}

func (b *BaseAdapter) ToJson() map[string]any {
	return map[string]any{
		"id":    b.id,
		"state": b.State(),
		"name":  b.name,
		"devices": utils.MapEntries(b.devices, func(name string, rd RegisteredDevice) map[string]any {
			return rd.Device.ToJson()
		}),
	}
}

func (b *BaseAdapter) Stop() error {
	b.mu.Lock()
	if b.cancel == nil {
		b.mu.Unlock()
		return nil
	}
	cancel := b.cancel
	b.cancel = nil
	b.mu.Unlock()

	cancel()

	_ = b.hooks.OnStop()

	return nil
}

func (b *BaseAdapter) RegisterDevice(d Device) error {
	b.mu.RLock()
	if _, ok := b.devices[d.Addr()]; ok {
		b.mu.RUnlock()
		return errors.New("device already registered")
	}
	b.mu.RUnlock()

	meta, cancel, err := b.hooks.OnRegisterDevice(d)
	if err != nil {
		return err
	}

	// store registration
	b.mu.Lock()
	b.devices[d.Addr()] = RegisteredDevice{
		Device: d,
		Meta:   meta,
		Cancel: cancel,
	}
	b.needRestart = true
	b.mu.Unlock()

	return nil
}

func (b *BaseAdapter) UnregisterDevice(addr string) error {
	b.mu.Lock()
	rd, ok := b.devices[addr]
	if !ok {
		b.mu.Unlock()
		return errors.New("device not found")
	}
	delete(b.devices, addr)
	b.mu.Unlock()

	if rd.Cancel != nil {
		rd.Cancel()
	}
	return b.hooks.OnUnregisterDevice(rd.Meta)
}

func (b *BaseAdapter) Context() context.Context {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.ctx == nil {
		return context.Background()
	}
	return b.ctx
}

func (b *BaseAdapter) DeviceMetas() []any {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]any, 0, len(b.devices))
	for _, rd := range b.devices {
		out = append(out, rd.Meta)
	}
	return out
}
