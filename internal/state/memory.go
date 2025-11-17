package state

import (
	"gohome/internal/core"
	"sync"
)

type MemoryStore struct {
	mu       sync.RWMutex
	adapters map[string]core.Adapter
	parsers  map[string]core.Parser
	devices  map[string]core.Device
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		adapters: make(map[string]core.Adapter),
		parsers:  make(map[string]core.Parser),
		devices:  make(map[string]core.Device),
	}
}

func (s *MemoryStore) GetAdapters() map[string]core.Adapter { return s.adapters }
func (s *MemoryStore) GetDevices() map[string]core.Device   { return s.devices }
func (s *MemoryStore) GetParsers() map[string]core.Parser   { return s.parsers }

func (s *MemoryStore) GetDevice(addr string) (core.Device, bool) {
	d, ok := s.devices[addr]
	return d, ok
}
func (s *MemoryStore) GetAdapter(id string) (core.Adapter, bool) {
	a, ok := s.adapters[id]
	return a, ok
}

func (s *MemoryStore) SaveDevice(d core.Device) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.devices[d.Addr()] = d
}

func (s *MemoryStore) SaveAdapter(a core.Adapter) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.adapters[a.Id()] = a
}

func (s *MemoryStore) SaveParser(p core.Parser) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.parsers[p.Name()] = p
}

func (s *MemoryStore) DeleteDevice(addr string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.devices, addr)
}
