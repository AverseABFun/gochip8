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

type StackRenderer interface {
	Parent() RawRenderer
	SetParent(rr RawRenderer)
	CanUseCurrentRawRenderer() bool
}
