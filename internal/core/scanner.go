package core

import (
	"context"
	"fmt"

	"sync"
)

type Scanner interface {
	Start(ctx context.Context) error
	Stop() error
	Out(addr string) <-chan Advertisment
	Started() bool
	AddAddress(address string)
}

type ScannerHooks interface {
	OnStart(ctx context.Context) error
	OnStop() error
}

type BaseScanner struct {
	ChOut     map[string]chan Advertisment
	mu        sync.Mutex
	Addresses []string
	started   bool
	hooks     ScannerHooks
}

func NewBaseScanner(addresses []string, hooks ScannerHooks) *BaseScanner {
	return &BaseScanner{
		ChOut:     make(map[string]chan Advertisment, 0),
		started:   false,
		Addresses: addresses,
		hooks:     hooks,
	}
}

func (s *BaseScanner) Started() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}

func (s *BaseScanner) AddAddress(address string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Addresses = append(s.Addresses, address)
}

func (s *BaseScanner) Stop() error {
	if !s.started {
		return nil
	}

	err := s.hooks.OnStop()

	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = false
	for _, v := range s.ChOut {
		close(v)
	}
	return err
}

func (s *BaseScanner) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.started {
		return fmt.Errorf("already started")
	}
	s.started = true
	for _, addr := range s.Addresses {
		ch := make(chan Advertisment, 20)
		s.ChOut[addr] = ch
	}

	return s.hooks.OnStart(ctx)
}

func (s *BaseScanner) Out(addr string) <-chan Advertisment {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ChOut[addr]
}
