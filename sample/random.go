package sample

import (
	"grpc/psm"
	"math/rand"
)

//KeyboardRandomizer

func randomKeyboardLayout() psm.KeyboardType_KeyboardLayout {

	switch rand.Intn(3) {
	case 1:
		return psm.KeyboardType_QWERTY
	case 2:
		return psm.KeyboardType_AZERTY
	default:
		return psm.KeyboardType_QWERTZ
	}
}

func randomKeyboard() bool {
	return rand.Intn(2) == 1 // return true if rand is one
}

// Screen_Panel randomizer

func randomPanel() psm.Screen_Panel {
	panel := psm.Screen_Panel(1 + rand.Intn(3))
	return panel
}

func randomScreenSize() float32 {
	return float32(13 + rand.Intn(50))
}

func randomScreenResolusion() *psm.Screen_Resolution {
	return &psm.Screen_Resolution{
		Width:  uint32(1920 + rand.Intn(3840)),
		Height: uint32(1080 + rand.Intn(2160)),
	}
}

func randomScreenMultiTouch() bool {
	return rand.Intn(2) == 1
}

// Memory Randomizer
