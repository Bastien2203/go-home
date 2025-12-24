package main

import (
	"hash/fnv"
	"log"
	"sync"

	"github.com/Bastien2203/go-home/utils"
	"github.com/brutella/hap/accessory"
	"github.com/brutella/hap/characteristic"
	"github.com/brutella/hap/service"
)

type HomekitManager struct {
	mu          sync.Mutex
	accessories map[string]*accessory.A
	server      *HomekitServer
}

func NewHomekitManager(server *HomekitServer) *HomekitManager {
	return &HomekitManager{
		accessories: make(map[string]*accessory.A),
		server:      server,
	}
}

func (s *HomekitManager) CreateAccessory(name string, id string, accType byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	accInfo := accessory.Info{
		Name:         name,
		SerialNumber: id,
		Manufacturer: "GoHome",
		Model:        "Virtual Device",
		Firmware:     "1.0.0",
	}

	acc := accessory.New(accInfo, accType)
	acc.Id = hashId(id)

	s.accessories[id] = acc
	log.Printf("[HomeKit] Device registered : %s", id)
	s.server.ScheduleReload(s.accessories)
}

func (s *HomekitManager) UpdateAccessory(id string, service *service.S) {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc := s.accessories[id]
	acc.AddS(service)
	log.Printf("[HomeKit] Device updated : %s", acc.Info.SerialNumber.Value())
	s.server.ScheduleReload(s.accessories)
}

func (s *HomekitManager) GetService(id string, serviceType string) *service.S {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc, exists := s.accessories[id]
	if !exists {
		return nil
	}
	return findService(acc, serviceType)
}

func (s *HomekitManager) CreateService(id string, service *service.S) *service.S {
	s.mu.Lock()
	defer s.mu.Unlock()

	acc, exists := s.accessories[id]
	if !exists {
		log.Printf("accessory %s doesnt exists", id)
		return nil
	}
	acc.AddS(service)
	s.server.ScheduleReload(s.accessories)
	return service
}

func (s *HomekitManager) AccessoryExists(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.accessories[id]
	return exists
}

func (s *HomekitManager) RemoveAccessory(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.accessories[id]; !ok {
		return
	}
	delete(s.accessories, id)
	log.Printf("[HomeKit] Device unregistered : %s", id)
	s.server.ScheduleReload(s.accessories)
}

func (s *HomekitManager) UdateCharacteristic(c *characteristic.C, val any) bool {
	if c == nil {
		return false
	}

	switch c.Format {
	case characteristic.FormatFloat:
		if v, ok := utils.ToFloat(val); ok {
			(&characteristic.Float{C: c}).SetValue(v)
		}
	case characteristic.FormatInt32, characteristic.FormatUInt8, characteristic.FormatUInt16:
		if v, ok := utils.ToInt(val); ok {
			(&characteristic.Int{C: c}).SetValue(v)
		}
	case characteristic.FormatBool:
		if v, ok := val.(bool); ok {
			(&characteristic.Bool{C: c}).SetValue(v)
		}
	case characteristic.FormatString:
		if v, ok := val.(string); ok {
			(&characteristic.String{C: c}).SetValue(v)
		}
	default:
		return false
	}
	return true
}

func findService(acc *accessory.A, serviceType string) *service.S {
	for _, s := range acc.Ss {
		if s.Type == serviceType {
			return s
		}
	}
	return nil
}

func hashId(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	val := h.Sum64()
	// Id 1 is reserved for bridge
	if val <= 1 {
		return val + 2
	}
	return val
}
