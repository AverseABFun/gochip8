package lib

import (
	"encoding/binary"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/averseabfun/gochip8/engine/impl"
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
	case 0x0000:
		break
	case 0x00E0:
		logging.Println(logging.MsgDebug, "Clearing screen")
		for x := range data.Memory.Display {
			for y := range data.Memory.Display[x] {
				data.Memory.Display[x][y] = 0
			}
		}
		data.Backend.FillBack(data.BackgroundColor)
	case 0x00EE:
		data.Registers.SP -= 1
		data.Registers.PC = data.Memory.Stack[data.Registers.SP] - 2
	default:
		switch inst & 0xF000 {
		case 0x1000:
			data.Registers.PC = (inst & 0xFFF) - 2
			break TopSwitch
		case 0x2000:
			data.Memory.Stack[data.Registers.SP] = data.Registers.PC + 2
			data.Registers.SP += 1
			data.Registers.PC = (inst & 0xFFF) - 2
			break TopSwitch
		case 0x3000:
			var reg = data.Registers.V[(inst&0xF00)>>8]
			var val = uint8(inst & 0xFF)
			if reg == val {
				data.Registers.PC += 2
			}
			break TopSwitch
		case 0x4000:
			var reg = data.Registers.V[(inst&0xF00)>>8]
			var val = uint8(inst & 0xFF)
			if reg != val {
				data.Registers.PC += 2
			}
			break TopSwitch
		case 0x5000:
			var reg = data.Registers.V[(inst&0xF00)>>8]
			var val = data.Registers.V[(inst&0xF0)>>4]
			if reg == val {
				data.Registers.PC += 2
			}
			break TopSwitch
		case 0x6000:
			var register = (inst & 0xF00) >> 8
			var storeData = uint8(inst & 0xFF)
			data.Registers.V[register] = storeData
			break TopSwitch
		case 0x7000:
			var register = (inst & 0xF00) >> 8
			var storeData = uint8(inst & 0xFF)
			data.Registers.V[register] += storeData
			break TopSwitch
		case 0x8000:
			switch inst & 0xF {
			case 0x0:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register1] = data.Registers.V[register2]
				break TopSwitch
			case 0x1:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register1] |= data.Registers.V[register2]
				data.Registers.V[0xF] = 0
				break TopSwitch
			case 0x2:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register1] &= data.Registers.V[register2]
				data.Registers.V[0xF] = 0
				break TopSwitch
			case 0x3:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register1] ^= data.Registers.V[register2]
				data.Registers.V[0xF] = 0
				break TopSwitch
			case 0x4:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				var result = data.Registers.V[register1] + data.Registers.V[register2]
				var overflow = result < data.Registers.V[register1] || result < data.Registers.V[register2]
				data.Registers.V[register1] = result
				if overflow {
					data.Registers.V[0xF] = 1
				} else {
					data.Registers.V[0xF] = 0
				}
				break TopSwitch
			case 0x5:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				var underflow = data.Registers.V[register2] > data.Registers.V[register1]
				data.Registers.V[register1] -= data.Registers.V[register2]
				if underflow {
					data.Registers.V[0xF] = 0
				} else {
					data.Registers.V[0xF] = 1
				}
				break TopSwitch
			case 0x6:
				var register = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register] = data.Registers.V[register2]
				var F = (data.Registers.V[register] & 1)
				data.Registers.V[register] >>= 1
				data.Registers.V[0xF] = F
				break TopSwitch
			case 0x7:
				var register1 = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				var underflow = data.Registers.V[register1] > data.Registers.V[register2]
				data.Registers.V[register1] = data.Registers.V[register2] - data.Registers.V[register1]
				if underflow {
					data.Registers.V[0xF] = 0
				} else {
					data.Registers.V[0xF] = 1
				}
				break TopSwitch
			case 0xE:
				var register = (inst & 0xF00) >> 8
				var register2 = (inst & 0xF0) >> 4
				data.Registers.V[register] = data.Registers.V[register2]
				var F = (data.Registers.V[register] & 0x80) >> 7
				data.Registers.V[register] <<= 1
				data.Registers.V[0xF] = F
				break TopSwitch
			default:
				logging.Panicf("got unknown instruction 0x%X at PC:0x%X", inst, data.Registers.PC)
			}
		case 0x9000:
			var reg = data.Registers.V[(inst&0xF00)>>8]
			var val = data.Registers.V[(inst&0xF0)>>4]
			if reg != val {
				data.Registers.PC += 2
			}
			break TopSwitch
		case 0xA000:
			data.Registers.I = inst & 0xFFF
			break TopSwitch
		case 0xB000:
			data.Registers.PC = (inst & 0xFFF) + uint16(data.Registers.V[0]) - 2
			break TopSwitch
		case 0xC000:
			var reg = data.Registers.V[(inst&0xF00)>>8]
			var val = uint8(inst & 0xFF)
			data.Registers.V[reg] = uint8(rand.Uint32()) & val
			break TopSwitch
		case 0xD000:
			var numberBytes = uint8(inst & 0xF)
			var x = data.Registers.V[(inst&0xF00)>>8]
			var y = data.Registers.V[(inst&0xF0)>>4]
			logging.Printf(logging.MsgDebug, "Writing %d bytes from 0x%X to screen pos (%d, %d)\n", numberBytes, data.Registers.I, x, y)
			data.Registers.V[0xF] = 0
			y %= 32
			x %= 64
			for offset := data.Registers.I; offset < data.Registers.I+uint16(numberBytes); offset++ {
				var b = data.Memory.AllMemory[offset]
				var bits = getBits(b)
				for i := 0; i < 8; i++ {
					var d uint8 = 0
					if bits[i] {
						d = 0xFF
					}
					if data.Memory.Display[x][y] != 0 {
						data.Registers.V[0xF] = 1
					}
					data.Memory.Display[x][y] ^= d
					if data.Memory.Display[x][y] > 0 {
						data.Backend.DrawBackPixel(uint32(x), uint32(y), data.ForegroundColor)
					} else {
						data.Backend.DrawBackPixel(uint32(x), uint32(y), data.BackgroundColor)
					}
					x++
				}
				x = x - 8
				y++
			}
			data.DoneProcessing = true
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
			default:
				logging.Panicf("got unknown instruction 0x%X at PC:0x%X", inst, data.Registers.PC)
			}
		case 0xF000:
			switch inst & 0xFF {
			case 0x07:
				var register = (inst & 0xF00) >> 8
				data.Registers.V[register] = data.Registers.DT
			case 0x0A:
				var register = (inst & 0xF00) >> 8
				var outputted = false
				f := func(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool) {
					logging.Println(logging.MsgDebug, "Got key")
					c := glfw.GetKeyName(key, scancode)
					var d uint8 = 0
					if action != glfw.Release {
						goto End
					}
					switch strings.ToLower(c) {
					case "1":
						d = 0x1
					case "2":
						d = 0x2
					case "3":
						d = 0x3
					case "4":
						d = 0xC
					case "q":
						d = 0x4
					case "w":
						d = 0x5
					case "e":
						d = 0x6
					case "r":
						d = 0xD
					case "a":
						d = 0x7
					case "s":
						d = 0x8
					case "d":
						d = 0x9
					case "f":
						d = 0xE
					case "z":
						d = 0xA
					case "x":
						d = 0x0
					case "c":
						d = 0xB
					case "v":
						d = 0xF
					default:
						goto End
					}
					data.Registers.V[register] = d
					logging.Println(logging.MsgDebug, "Got released key")
					outputted = true
				End:
					return false
				}
				var grab1, _ = data.Backend.PopGrabber()
				var grabber = data.Backend.PushGrabber(impl.FuncGrabber{Function: f})
				logging.Println(logging.MsgDebug, "Pushed")
				for !outputted {
				}
				data.Backend.PopGrabberAt(grabber)
				data.Backend.PushGrabber(grab1)
			case 0x15:
				var register = (inst & 0xF00) >> 8
				data.Registers.DT = data.Registers.V[register]
			case 0x18:
				var register = (inst & 0xF00) >> 8
				data.Registers.ST = data.Registers.V[register]
				var err error
				data.CurrentToneID, err = data.AudioBackend.PlayTone(500)
				if err != nil {
					panic(err)
				}
				data.Playing = true
			case 0x1E:
				var register = (inst & 0xF00) >> 8
				data.Registers.I += uint16(data.Registers.V[register])
			case 0x29:
				var register = (inst & 0xF00) >> 8
				data.Registers.I = (uint16(data.Registers.V[register]) * 5) + 50
			case 0x33:
				var register = (inst & 0xF00) >> 8
				var value = data.Registers.V[register]
				// Ones-place
				data.Memory.AllMemory[data.Registers.I+2] = value % 10
				value /= 10

				// Tens-place
				data.Memory.AllMemory[data.Registers.I+1] = value % 10
				value /= 10

				// Hundreds-place
				data.Memory.AllMemory[data.Registers.I] = value % 10
			case 0x55:
				var register = (inst & 0xF00) >> 8
				var offset uint16
				for offset = data.Registers.I; offset <= data.Registers.I+register; offset++ {
					data.Memory.AllMemory[offset] = data.Registers.V[offset-data.Registers.I]
				}
				data.Registers.I = offset
			case 0x65:
				var register = (inst & 0xF00) >> 8
				var offset uint16
				for offset = data.Registers.I; offset <= data.Registers.I+register; offset++ {
					data.Registers.V[offset-data.Registers.I] = data.Memory.AllMemory[offset]
				}
				data.Registers.I = offset
			default:
				logging.Panicf("got unknown instruction 0x%X at PC:0x%X", inst, data.Registers.PC)
			}
		default:
			logging.Panicf("got unknown instruction 0x%X at PC:0x%X", inst, data.Registers.PC)
		}
	}
	data.Registers.PC += 2
}

func (data *Chip8Data) TickAll() {
	f := func(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool) {
		c := glfw.GetKeyName(key, scancode)
		var d = 0
		if action == glfw.Repeat {
			goto End
		}
		switch strings.ToLower(c) {
		case "1":
			d = 0x1
		case "2":
			d = 0x2
		case "3":
			d = 0x3
		case "4":
			d = 0xC
		case "q":
			d = 0x4
		case "w":
			d = 0x5
		case "e":
			d = 0x6
		case "r":
			d = 0xD
		case "a":
			d = 0x7
		case "s":
			d = 0x8
		case "d":
			d = 0x9
		case "f":
			d = 0xE
		case "z":
			d = 0xA
		case "x":
			d = 0x0
		case "c":
			d = 0xB
		case "v":
			d = 0xF
		default:
			goto End
		}
		data.KeysPressed.Keys[d] = action == glfw.Press
		logging.Println(logging.MsgDebug, data.KeysPressed.Keys)
	End:
		return false
	}
	var grabber = data.Backend.PushGrabber(impl.FuncGrabber{Function: f})
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		if recover() == nil {
			return
		}
		done <- syscall.SIGINT
	}()

	go func() {
		var startTime time.Time
		var duration = time.Duration(16670000 * float64(time.Nanosecond))
		for {
			startTime = time.Now()
			(*data.DrawTimer) += 1
			time.Sleep(time.Duration(math.Max(float64(duration-time.Since(startTime)), 0)))
		}
	}()

	var startTime time.Time
	var duration = time.Duration(16670000 * float64(time.Nanosecond))
	for {
		startTime = time.Now()
		data.Backend.TickRenderer()
		(*data.DrawTimer) += 1
		select {
		case <-done:
			signal.Stop(done)
			data.Backend.PopGrabberAt(grabber)
			return
		default:
			if data.Registers.DT > 0 {
				data.Registers.DT -= 1
			}
			if data.Registers.ST > 0 {
				data.Registers.ST -= 1
			} else if data.Playing && data.CurrentToneID != 0 {
				data.AudioBackend.StopAll()
				data.CurrentToneID = 0
				data.Playing = false
			}
			for i := 0; i < data.InstPerFrame; i++ {
				data.TickSingle()
				if data.DoneProcessing {
					break
				}
			}
			data.DoneProcessing = false
			for time.Since(startTime) < duration {
			}
		}
	}
}
