package impl

import "github.com/go-gl/glfw/v3.3/glfw"

type FuncGrabber struct {
	Function func(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool)
}

func (fg FuncGrabber) GrabKey(key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) (continueSearching bool) {
	return fg.Function(key, scancode, action, mods)
}
