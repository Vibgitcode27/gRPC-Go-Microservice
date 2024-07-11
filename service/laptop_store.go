package service

import (
	"errors"
	"grpc/psm"
	"sync"

	"github.com/jinzhu/copier"
)

var ErrAlreadyExists = errors.New("laptop already exists")

type LaptopStore interface {
	Save(laptop *psm.Laptop) error
	// Find(id string) (*psm.Laptop, error)
	Search(filter *psm.Filter, found func(laptop *psm.Laptop) error) error
}

type InMemoryLaptopStore struct {
	mutex sync.RWMutex
	data  map[string]*psm.Laptop
}

// NewInMemoryLaptopStore returns a new InMemoryLaptopStore
func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*psm.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *psm.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return ErrAlreadyExists
	}

	// deep copy of laptop object
	other, err := deepCopy(laptop)
	if err != nil {
		return err
	}

	store.data[laptop.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) Find(id string) (*psm.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	// deep copy of laptop object
	other := &psm.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, err
	}

	return other, nil
}

func (store *InMemoryLaptopStore) Search(filter *psm.Filter, found func(laptop *psm.Laptop) error) error {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, laptop := range store.data {
		if isQualified(filter, laptop) {
			other, err := deepCopy(laptop)
			if err != nil {
				return err
			}
			err = found(other)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isQualified(filter *psm.Filter, laptop *psm.Laptop) bool {
	if filter.GetMaxPriceInr() > 0 && laptop.GetPriceInr() > filter.GetMaxPriceInr() {
		return false
	}

	if filter.GetMinCpuCores() > 0 && laptop.GetCpu().GetCores() < filter.GetMinCpuCores() {
		return false
	}

	if filter.GetMinCpuGhz() > 0 && laptop.GetCpu().GetMinGhz() < filter.GetMinCpuGhz() {
		return false
	}

	if toBit(filter.GetRam()) > toBit(laptop.GetRam()) {
		return false
	}

	return true
}

func toBit(ram *psm.Memory) uint64 {
	value := ram.GetValue()
	unit := ram.GetUnit()

	switch unit {
	case psm.Memory_BIT:
		return value
	case psm.Memory_BYTE:
		return value << 3
	case psm.Memory_KILOBYTE:
		return value << 13
	case psm.Memory_MEGABYTE:
		return value << 23
	case psm.Memory_GIGABYTE:
		return value << 33
	case psm.Memory_TERABYTE:
		return value << 43
	default:
		return 0
	}
}

func deepCopy(laptop *psm.Laptop) (*psm.Laptop, error) {
	other := &psm.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, err
	}
	return other, nil
}
