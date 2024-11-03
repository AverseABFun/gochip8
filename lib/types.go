package lib

import (
	"github.com/averseabfun/gochip8/engine/interfaces"
	"github.com/averseabfun/gochip8/engine/types"
)

type Chip8Data struct {
	Memory          Memory
	Registers       Registers
	KeysPressed     KeysPressed
	Backend         interfaces.FullIO
	AudioBackend    interfaces.AudioRenderer
	Initialized     bool
	ClockSpeed      float64
	CurrentToneID   interfaces.AudioID
	Playing         bool
	DrawTimer       *int
	busyWork        bool
	BackgroundColor types.Color
	ForegroundColor types.Color
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
	data.Backend.FillBack(data.BackgroundColor)
	data.Backend.TickRenderer()
	if err := data.AudioBackend.InitAudio(); err != nil {
		panic(err)
	}
	if data.ClockSpeed == 0 {
		data.ClockSpeed = 500
	}
	var val = 0
	data.DrawTimer = &val
	data.Initialized = true
}
