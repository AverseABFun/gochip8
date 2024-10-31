package lib

import (
	"encoding/binary"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/averseabfun/gochip8/engine/impl"
	"github.com/averseabfun/gochip8/engine/types"
	"github.com/averseabfun/gochip8/logging"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func getBits(b byte) [8]bool {
	return [8]bool{
		(b>>7)&1 == 1,
		(b>>6)&1 == 1,
		(b>>5)&1 == 1,
		(b>>4)&1 == 1,
		(b>>3)&1 == 1,
		(b>>2)&1 == 1,
		(b>>1)&1 == 1,
		(b>>0)&1 == 1,
	}
}

func (data *Chip8Data) TickSingle() {
	data.InitalizeData()
	var bytes = []byte{data.Memory.AllMemory[data.Registers.PC], data.Memory.AllMemory[data.Registers.PC+1]}
	var inst = binary.BigEndian.Uint16(bytes)

TopSwitch:
	switch inst {
	case 0x00E0:
		logging.Println(logging.MsgDebug, "Clearing screen")
		for x := range data.Memory.Display {
			for y := range data.Memory.Display[x] {
				data.Memory.Display[x][y] = 0
			}
		}
		data.Backend.FillBack(types.FromRGBNoErr(0, 0, 0))
		data.Backend.TickRenderer()
	case 0x00EE:
		data.Registers.PC = data.Memory.Stack[data.Registers.SP] - 2
		data.Registers.SP -= 1
	default:
		switch inst & 0xF000 {
		case 0x6000:
			var register = (inst & 0xF00) >> 8
			var storeData = uint8(inst & 0xFF)
			data.Registers.V[register] = storeData
			break TopSwitch
		case 0xA000:
			data.Registers.I = inst & 0xFFF
			break TopSwitch
		case 0x1000:
			data.Registers.PC = (inst & 0xFFF) - 2
			break TopSwitch
		case 0x2000:
			data.Registers.SP += 1
			data.Memory.Stack[data.Registers.SP] = data.Registers.PC
			data.Registers.PC = (inst & 0xFFF) - 2
			break TopSwitch
		case 0xD000:
			var numberBytes = uint8(inst & 0xF)
			var x = data.Registers.V[(inst&0xF00)>>8]
			var y = data.Registers.V[(inst&0xF0)>>4]
			logging.Printf(logging.MsgDebug, "Writing %d bytes from 0x%X to screen pos (%d, %d)\n", numberBytes, data.Registers.I, x, y)
			data.Registers.V[0xF] = 0
			for offset := data.Registers.I; offset < data.Registers.I+uint16(numberBytes); offset++ {
				y %= 64
				var b = data.Memory.AllMemory[offset]
				var bits = getBits(b)
				for i := 0; i < 8; i++ {
					x %= 128
					var d uint8 = 0
					if bits[i] {
						d = 0xFF
					}
					if data.Memory.Display[x][y] != 0 {
						data.Registers.V[0xF] = 1
					}
					data.Memory.Display[x][y] ^= d
					if data.Memory.Display[x][y] > 0 {
						data.Backend.DrawBackPixel(uint32(x), uint32(y), types.FromRGBNoErr(types.MAX_UINT6, types.MAX_UINT6, types.MAX_UINT6))
					} else {
						data.Backend.DrawBackPixel(uint32(x), uint32(y), types.FromRGBNoErr(0, 0, 0))
					}
					x++
				}
				x = x - 8
				y++
			}
			data.Backend.TickRenderer()
			break TopSwitch
		case 0x7000:
			var register = (inst & 0xF00) >> 8
			var storeData = uint8(inst & 0xFF)
			data.Registers.V[register] += storeData
			break TopSwitch
		case 0x4000:
			var register = (inst & 0xF00) >> 8
			var storeData = uint8(inst & 0xFF)
			if data.Registers.V[register] != storeData {
				data.Registers.PC += 2
			}
			break TopSwitch
		case 0xE000:
			switch inst & 0xFF {
			case 0x9E:
				var register = (inst & 0xF00) >> 8
				if data.KeysPressed.Keys[data.Registers.V[register]] {
					data.Registers.PC += 2
				}
				break TopSwitch
			case 0xA1:
				var register = (inst & 0xF00) >> 8
				if !data.KeysPressed.Keys[data.Registers.V[register]] {
					data.Registers.PC += 2
				}
				break TopSwitch
			}
		default:
			logging.Panicf("got unknown instruction 0x%X at PC:0x%X", inst, data.Registers.PC)
		}
	}
	data.Registers.PC += 2
}

func (data *Chip8Data) TickAll() {
	ch := make(chan string)
	f := func(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool) {
		ch <- glfw.GetKeyName(key, scancode)
		return false
	}
	data.Backend.PushGrabber(impl.FuncGrabber{Function: f})
	var d = 0
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-done:
			signal.Stop(done)
			return
		case c := <-ch:
			if !strings.HasSuffix(c, ";") {
				switch strings.ToLower(c)[0] {
				case '1':
					d = 0x1
				case '2':
					d = 0x2
				case '3':
					d = 0x3
				case '4':
					d = 0xC
				case 'q':
					d = 0x4
				case 'w':
					d = 0x5
				case 'e':
					d = 0x6
				case 'r':
					d = 0xD
				case 'a':
					d = 0x7
				case 's':
					d = 0x8
				case 'd':
					d = 0x9
				case 'f':
					d = 0xE
				case 'z':
					d = 0xA
				case 'x':
					d = 0x0
				case 'c':
					d = 0xB
				case 'v':
					d = 0xF
				default:
					goto EndIf
				}
				data.KeysPressed.Keys = [16]bool{}
				data.KeysPressed.Keys[d] = true
				logging.Println(logging.MsgDebug, data.KeysPressed.Keys)
			}
		EndIf:
			break
		default:
			data.TickSingle()
		}
	}
}
