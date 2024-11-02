package interfaces

import "github.com/averseabfun/gochip8/engine/types"

type RawRenderer interface {
	InitRenderer(windowName string, width uint32, height uint32) error
	GetSize() types.Point
	TickRenderer()
	ShouldQuit() bool
	DeinitRenderer() error
	DrawBackPixel(x uint32, y uint32, color types.Color) error
	FillBack(color types.Color) error
}

type AudioRenderer interface {
	InitAudio() error
	DeinitAudio() error
	PlayTone(freqHz float64) (AudioID, error)
	Playing(tone AudioID) bool
	PauseTone(tone AudioID) error
	ResumeTone(tone AudioID) error
	StopTone(tone AudioID) error
	StopAll() error
}

type AudioID uint64

type StackRenderer interface {
	Parent() RawRenderer
	SetParent(rr RawRenderer)
	CanUseCurrentRawRenderer() bool
}
