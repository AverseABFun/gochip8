package lib

import (
	"io"
	"unsafe"

	"github.com/averseabfun/gochip8/logging"
)

type Memory struct {
	Font                *[80]uint8
	InterpreterReserved *[0x1FF]uint8
	MainMemory          *[0xE00]uint8
	AllMemory           [0xFFF]uint8
	Stack               [16]uint16
	Display             [128][64]uint8
}

func CreateEmptyMemory() Memory {
	var out = Memory{}
	out.InterpreterReserved = (*[0x1FF]uint8)(unsafe.Pointer(&out.AllMemory))
	out.MainMemory = (*[0xE00]uint8)(unsafe.Pointer(&out.AllMemory[0x200]))
	out.Font = (*[80]uint8)(unsafe.Pointer(&(*out.InterpreterReserved)[50]))
	copy(out.Font[:], Font[:])
	logging.Println(logging.MsgDebug, "Created empty memory")
	return out
}

func (m *Memory) LoadMemory(i io.Reader) error {
	var err error = nil
	var offset = 0
	for err == nil {
		var bytes = make([]byte, 1)
		_, err = i.Read(bytes)
		if err != nil {
			logging.Printf(logging.MsgInfo, "Read %d bytes\n", offset)
			break
		}
		m.AllMemory[offset+0x200] = bytes[0]
		offset++
	}
	return err
}

type Registers struct {
	V  [16]uint8
	I  uint16
	DT uint8
	ST uint8
	PC uint16
	SP uint8
}

type KeysPressed struct {
	Keys [16]bool
}
