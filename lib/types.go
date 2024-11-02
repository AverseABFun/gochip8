package lib

import (
	"github.com/averseabfun/gochip8/engine/interfaces"
)

type Chip8Data struct {
	Memory        Memory
	Registers     Registers
	KeysPressed   KeysPressed
	Backend       interfaces.FullIO
	AudioBackend  interfaces.AudioRenderer
	Initialized   bool
	ClockSpeed    float64
	CurrentToneID interfaces.AudioID
	Playing       bool
}

func (data *Chip8Data) InitalizeData() {
	if data.Initialized {
		return
	}
	data.Registers.PC = 0x200
	if err := data.Backend.InitRenderer("GoChip8", 128, 64); err != nil {
		panic(err)
	}
	data.Backend.TickRenderer()
	if err := data.AudioBackend.InitAudio(); err != nil {
		panic(err)
	}
	data.ClockSpeed = 500
	data.Initialized = true
}
