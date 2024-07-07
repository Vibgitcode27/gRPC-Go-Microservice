package sample

import (
	"grpc/psm"
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
		Brand: "Intel",
		Cores: 6,
	}
	return cpu
}

func Display() *psm.Screen {
	return &psm.Screen{
		Panel:      randomPanel(),
		SizeInch:   randomScreenSize(),
		Resolution: randomScreenResolusion(),
		MultiTouch: randomScreenMultiTouch(),
	}
}
