package lib

import (
	"github.com/averseabfun/gochip8/engine/interfaces"
)

type Chip8Data struct {
	Memory      Memory
	Registers   Registers
	KeysPressed KeysPressed
	Backend     interfaces.FullIO
	Initialized bool
	ClockSpeed  float64
}

func (data *Chip8Data) InitalizeData() {
	if data.Initialized {
		return
	}
	data.Registers.PC = 0x200
	data.Backend.InitRenderer("GoChip8", 128, 64)
	data.Backend.TickRenderer()
	data.ClockSpeed = 500
	data.Initialized = true
}
