package interfaces

import "github.com/go-gl/glfw/v3.3/glfw"

type KeyGrabber interface {
	GrabKey(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool)
}

type MouseGrabber interface {
	GrabMouse(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey, posX float64, posY float64) (continueSearching bool)
}

type KeyProvider interface {
	PushGrabber(grabber KeyGrabber) (index uint32)
	PopGrabber() (KeyGrabber, error)
	PushGrabberAt(grabber KeyGrabber, index uint32)
	PopGrabberAt(index uint32) (KeyGrabber, error)
}

type MouseProvider interface {
	PushMouseGrabber(grabber MouseGrabber)
	PopMouseGrabber() (MouseGrabber, error)
	PushMouseGrabberAt(grabber MouseGrabber, index uint32)
	PopMouseGrabberAt(index uint32) (MouseGrabber, error)
}

type FullIO interface {
	KeyProvider
	RawRenderer
}
