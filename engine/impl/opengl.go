package impl

import (
	"errors"
	"slices"

	"github.com/averseabfun/gochip8/engine/interfaces"
	"github.com/averseabfun/gochip8/engine/types"
	"github.com/averseabfun/gochip8/logging"
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type OpenGL struct {
	window        *glfw.Window
	pixels        [][]types.Color
	texture       uint32
	width         uint32
	height        uint32
	shouldClose   bool
	grabbers      []interfaces.KeyGrabber
	mouseGrabbers []interfaces.MouseGrabber
	focused       bool
}

func (rr *OpenGL) InitRenderer(windowName string, width uint32, height uint32) error {
	if err := glfw.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 5)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.False)

	window, err := glfw.CreateWindow(int(width)*8, int(height)*8, windowName, nil, nil)
	if err != nil {
		glfw.Terminate()
		rr.shouldClose = true
		return err
	}
	rr.window = window

	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	if err := gl.Init(); err != nil {
		return err
	}

	gl.Viewport(0, 0, int32(width)*8, int32(height)*8)

	rr.pixels = make([][]types.Color, 0)
	for x := uint32(0); x < width; x++ {
		var new = make([]types.Color, height)
		for y := uint32(0); y < height; y++ {
			new = append(new, types.FromRGBNoErr(0, 0, 0))
		}
		rr.pixels = append(rr.pixels, new)
	}

	gl.GenTextures(1, &rr.texture)
	gl.BindTexture(gl.TEXTURE_2D, rr.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	rr.width = width
	rr.height = height
	rr.shouldClose = false

	window.SetKeyCallback(rr.key_callback)
	window.SetMouseButtonCallback(rr.mouse_button_callback)
	window.SetFocusCallback(rr.focus_callback)
	rr.focused = true

	return nil
}

func (rr *OpenGL) GetSize() types.Point {
	var x, y = rr.window.GetSize()
	return types.Point{X: float64(x) / 4, Y: float64(y) / 4}
}

func (rr *OpenGL) DeinitRenderer() error {
	gl.DeleteTextures(1, &rr.texture)
	rr.window.Destroy()
	glfw.Terminate()
	rr.shouldClose = true
	return nil
}

func (rr *OpenGL) ShouldQuit() bool {
	return rr.shouldClose
}

func rotate90(matrix [][]types.Color) [][]types.Color {
	if len(matrix) == 0 {
		return matrix
	}

	// Get dimensions of the matrix
	rows := len(matrix)
	cols := len(matrix[0])

	// Create a new matrix to store the rotated result
	rotated := make([][]types.Color, cols)
	for i := range rotated {
		rotated[i] = make([]types.Color, rows)
	}

	// Rotate the matrix by 90 degrees clockwise
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			rotated[c][rows-1-r] = matrix[r][c]
		}
	}

	return rotated
}

func (rr *OpenGL) GetRGBArray() []uint8 {
	var rotated = rotate90(rr.pixels)

	var out = make([]uint8, 0, rr.width*rr.height*3)
	for y := 0; y < int(rr.height); y++ {
		for x := 0; x < int(rr.width); x++ {
			idx := rotated[y][x]
			out = append(out, uint8(idx.R)*4, uint8(idx.G)*4, uint8(idx.B)*4)
		}
	}
	return out
}

func (rr *OpenGL) TickRenderer() {
	if rr.window.ShouldClose() {
		rr.DeinitRenderer()
		return
	}
	glfw.PollEvents()
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Bind and update the texture
	gl.Enable(gl.TEXTURE_2D)
	gl.BindTexture(gl.TEXTURE_2D, rr.texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.SRGB, int32(rr.width), int32(rr.height), 0, gl.RGB, gl.UNSIGNED_BYTE, gl.Ptr(rr.GetRGBArray()))

	// Ensure to unbind the texture
	gl.BindTexture(gl.TEXTURE_2D, 0)

	// Create and bind framebuffer for reading
	var readFboId uint32
	gl.GenFramebuffers(1, &readFboId)
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, readFboId)
	gl.FramebufferTexture2D(gl.READ_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, rr.texture, 0)

	gl.BlitFramebuffer(
		0, 0, int32(rr.width), int32(rr.height),
		0, 0, int32(rr.width)*8, int32(rr.height)*8,
		gl.COLOR_BUFFER_BIT, gl.NEAREST,
	)

	// Unbind framebuffer and delete
	gl.BindFramebuffer(gl.READ_FRAMEBUFFER, 0)
	gl.DeleteFramebuffers(1, &readFboId)
	gl.Disable(gl.TEXTURE_2D)

	gl.Flush()
}

func (rr *OpenGL) DrawBackPixel(x uint32, y uint32, color types.Color) error {
	if x >= rr.width {
		return errors.New("got x over the width of the window")
	}
	if y >= rr.height {
		return errors.New("got y over the height of the window")
	}
	rr.pixels[rr.width-1-x][rr.height-1-y] = color
	logging.Printf(logging.MsgSpam, "Writing color (%d, %d, %d) to location (%d, %d)\n", color.R, color.G, color.B, x, y)
	return nil
}

func (rr *OpenGL) FillBack(color types.Color) error {
	for x := uint32(0); x < rr.width; x++ {
		for y := uint32(0); y < rr.height; y++ {
			rr.pixels[x][y] = color
		}
	}
	return nil
}

func (rr *OpenGL) key_callback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if !rr.focused {
		return
	}
	for i, grabber := range rr.grabbers {
		logging.Printf(logging.MsgDebug, "Calling grabber %d/%d\n", i, len(rr.grabbers))
		if grabber.GrabKey(key, scancode, action, mods) {
			break
		}
		logging.Println(logging.MsgDebug, "Moving on")
	}
}

func (rr *OpenGL) mouse_button_callback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if !rr.focused {
		return
	}
	var posX, posY = rr.window.GetCursorPos()
	posX /= 4
	posY /= 4
	for _, grabber := range rr.mouseGrabbers {
		if grabber.GrabMouse(button, action, mods, posX, posY) {
			break
		}
	}
}

func (rr *OpenGL) focus_callback(w *glfw.Window, focused bool) {
	rr.focused = focused
}

func (rr *OpenGL) PushGrabber(grabber interfaces.KeyGrabber) (index uint32) {
	rr.grabbers = append(rr.grabbers, grabber)
	return uint32(len(rr.grabbers)) - 1
}

func (rr *OpenGL) PopGrabber() (interfaces.KeyGrabber, error) {
	if len(rr.grabbers) == 0 {
		return nil, errors.New("empty stack")
	}
	var out = rr.grabbers[len(rr.grabbers)-1]
	rr.grabbers = slices.Delete(rr.grabbers, len(rr.grabbers)-1, len(rr.grabbers))
	return out, nil
}

func (rr *OpenGL) PushGrabberAt(grabber interfaces.KeyGrabber, index uint32) {
	rr.grabbers = slices.Insert(rr.grabbers, int(index), grabber)
}

func (rr *OpenGL) PopGrabberAt(index uint32) (interfaces.KeyGrabber, error) {
	if int(index) >= len(rr.grabbers) {
		return nil, errors.New("too small stack")
	}
	var out = rr.grabbers[index]
	rr.grabbers = slices.Delete(rr.grabbers, int(index), int(index)+1)
	return out, nil
}

func (rr *OpenGL) PushMouseGrabber(grabber interfaces.MouseGrabber) {
	rr.mouseGrabbers = append(rr.mouseGrabbers, grabber)
}

func (rr *OpenGL) PopMouseGrabber() (interfaces.MouseGrabber, error) {
	if len(rr.mouseGrabbers) == 0 {
		return nil, errors.New("empty stack")
	}
	var out = rr.mouseGrabbers[len(rr.grabbers)-1]
	rr.mouseGrabbers = slices.Delete(rr.mouseGrabbers, len(rr.mouseGrabbers)-1, len(rr.mouseGrabbers))
	return out, nil
}

func (rr *OpenGL) PushMouseGrabberAt(grabber interfaces.MouseGrabber, index uint32) {
	rr.mouseGrabbers = slices.Insert(rr.mouseGrabbers, int(index), grabber)
}

func (rr *OpenGL) PopMouseGrabberAt(index uint32) (interfaces.MouseGrabber, error) {
	if int(index) >= len(rr.mouseGrabbers) {
		return nil, errors.New("too small stack")
	}
	var out = rr.mouseGrabbers[index]
	rr.mouseGrabbers = slices.Delete(rr.mouseGrabbers, int(index), int(index)+1)
	return out, nil
}
