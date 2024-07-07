package sample

import (
	"grpc/psm"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewKeyboard() *psm.KeyboardType {
	keyBoard := &psm.KeyboardType{
		Layout:  randomKeyboardLayout(),
		Backlit: randomKeyboard(),
	}
	return keyBoard
}

func CPU() *psm.CPU {
	cpu := &psm.CPU{
		Brand:  "Intel",
		Name:   "Core i7",
		Cores:  6,
		MinGhz: 2.6,
		MaxGhz: 4.5,
	}
	return cpu
}

func GPU() *psm.GPU {
	return &psm.GPU{
		Brand:  "NVIDIA",
		Name:   "GeForce GTX 1080",
		Memory: randomMemory(),
		MinGhz: 1.6,
		MaxGhz: 1.8,
	}
}

func Display() *psm.Screen {
	return &psm.Screen{
		Panel:      randomPanel(),
		SizeInch:   randomScreenSize(),
		Resolution: randomScreenResolusion(),
		MultiTouch: randomScreenMultiTouch(),
	}
}

func randomStorage() *psm.Storage {
	return &psm.Storage{
		Driver: randomStorageDriver(),
		Memory: randomMemory(),
	}
}

func Laptop() *psm.Laptop {
	return &psm.Laptop{
		Id:          randomUUID(),
		Brand:       "Lenovo",
		Name:        "Thinkpad X1 Carbon",
		Cpu:         CPU(),
		Gpu:         []*psm.GPU{GPU(), GPU()},
		Ram:         randomMemory(),
		Storages:    []*psm.Storage{randomStorage(), randomStorage()},
		Keyboard:    NewKeyboard(),
		Screen:      Display(),
		Weight:      &psm.Laptop_WeightKg{WeightKg: 1.25},
		ReleaseYear: 2004,
		PriceInr:    200000,
		UpdatedAt:   timestamppb.Now(),
	}
}
