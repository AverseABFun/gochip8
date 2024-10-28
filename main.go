package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/averseabfun/gochip8/engine/impl"
)

func main() {
	var data = Chip8Data{
		Memory:      CreateEmptyMemory(),
		Registers:   Registers{},
		KeysPressed: KeysPressed{},
		Backend:     &impl.OpenGL{},
	}
	var file, err = os.Open("Cave.ch8")
	if err != nil {
		panic(err)
	}
	data.Memory.LoadMemory(file)
	data.InitalizeData()
	data.TickAll()
	var c = ""
	var d = 0
	for {
		fmt.Scanf("%s", &c)
		if !strings.HasSuffix(c, ";") && c != "" {
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
			fmt.Println(data.KeysPressed.Keys)
		} else {
			c = strings.TrimSuffix(c, ";")
			switch strings.ToLower(c) {
			case "step", "s", "st", "ste":
				fmt.Println("Stepping...")
				data.TickSingle()
			case "continue", "continu", "contin", "conti", "cont", "con", "co", "c":
				data.TickAll()
			case "screen", "scree", "scre", "scr", "sc", "r":
				for _, val := range data.Memory.Display[:] {
					fmt.Printf("%v\n", val)
				}
			case "quit", "qui", "qu", "q":
				os.Exit(0)
			case "dump", "dum", "du", "d":
				fmt.Printf(
					"Registers: \n\tV0: 0x%X\n\tV1: 0x%X\n\tV2: 0x%X\n\tV3: 0x%X\n\tV4: 0x%X\n\tV5: 0x%X\n\tV6: 0x%X\n\tV7: 0x%X\n\tV8: 0x%X\n\tV9: 0x%X\n\tVA: 0x%X\n\tVB: 0x%X\n\tVC: 0x%X\n\tVD: 0x%X\n\tVE: 0x%X\n\tVF: 0x%X\n",
					data.Registers.V[0x0],
					data.Registers.V[0x1],
					data.Registers.V[0x2],
					data.Registers.V[0x3],
					data.Registers.V[0x4],
					data.Registers.V[0x5],
					data.Registers.V[0x6],
					data.Registers.V[0x7],
					data.Registers.V[0x8],
					data.Registers.V[0x9],
					data.Registers.V[0xA],
					data.Registers.V[0xB],
					data.Registers.V[0xC],
					data.Registers.V[0xD],
					data.Registers.V[0xE],
					data.Registers.V[0xF],
				)
				fmt.Printf(
					"\tI: 0x%X\n\tDT: 0x%X\n\tST: 0x%X\n\tPC: 0x%X\n\tSP: 0x%X\n",
					data.Registers.I,
					data.Registers.DT,
					data.Registers.ST,
					data.Registers.PC,
					data.Registers.SP,
				)
			}
		}
	EndIf:
	}
}
