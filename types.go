package main

import (
	"github.com/averseabfun/gochip8/engine/interfaces"
	"github.com/averseabfun/gochip8/engine/types"
)

type Chip8Data struct {
	Memory      Memory
	Registers   Registers
	KeysPressed KeysPressed
	Backend     interfaces.FullIO
	Initialized bool
}

func (data *Chip8Data) InitalizeData() {
	if data.Initialized {
		return
	}
	data.Registers.PC = 0x200
	data.Backend.InitRenderer("GoChip8", 128, 64)
	data.Backend.DrawBackPixel(uint32(5), uint32(5), types.FromRGBNoErr(types.MAX_UINT6, types.MAX_UINT6, types.MAX_UINT6))
	data.Backend.TickRenderer()
	data.Initialized = true
}
