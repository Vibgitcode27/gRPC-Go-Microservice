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
	other := &psm.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return err
	}

	store.data[laptop.Id] = other
	return nil
}
