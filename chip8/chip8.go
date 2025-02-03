package chip8

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	scancodeHexMap map[sdl.Scancode]uint8 = map[sdl.Scancode]uint8{
		sdl.SCANCODE_1: 0x1, sdl.SCANCODE_2: 0x2, sdl.SCANCODE_3: 0x3, sdl.SCANCODE_4: 0xC,
		sdl.SCANCODE_Q: 0x4, sdl.SCANCODE_W: 0x5, sdl.SCANCODE_E: 0x6, sdl.SCANCODE_R: 0xD,
		sdl.SCANCODE_A: 0x7, sdl.SCANCODE_S: 0x8, sdl.SCANCODE_D: 0x9, sdl.SCANCODE_F: 0xE,
		sdl.SCANCODE_Z: 0xA, sdl.SCANCODE_X: 0x0, sdl.SCANCODE_C: 0xB, sdl.SCANCODE_V: 0xF,
	}
)

type Options struct {
	Scale int
	TPS   int
}

type Chip8 struct {
	opts         Options
	vm           *VirtualMachine
	keyBuffer    [NumKeys]bool
	lastKeyPress int
	screenBuffer [ScreenSize]bool
}

func New(opts Options) *Chip8 {
	interpreter := &Chip8{
		opts:         opts,
		keyBuffer:    [NumKeys]bool{},
		lastKeyPress: -1,
		screenBuffer: [ScreenSize]bool{},
	}

	interpreter.vm = NewVirtualMachine(interpreter, interpreter)

	return interpreter
}

func (c *Chip8) LoadROM(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return c.vm.LoadProgram(data...)
}

func (c *Chip8) Run() error {
	log.Info("Initializing CHIP-8")

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"CHIP-8",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(ScreenWidth*c.opts.Scale),
		int32(ScreenHeight*c.opts.Scale),
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN,
	)

	if err != nil {
		return err
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		return err
	}
	defer renderer.Destroy()

	rect := sdl.Rect{
		W: int32(c.opts.Scale),
		H: int32(c.opts.Scale),
	}

	log.Info("Starting CHIP-8")

	var (
		paused  = false
		running = true
	)

	// Runs at roughly 60 FPS
	for range time.Tick(time.Millisecond * 1000 / 60) {
		if !running {
			break
		}

		switch c.handleEvent() {
		case EventQuit:
			running = false

		case EventPause:
			paused = !paused

		case EventReset:
			c.vm.Reset()
			c.updateScreen(renderer, &rect)

		case EventNextCycle:
			if paused {
				c.vm.RunCycle()
				c.updateScreen(renderer, &rect)
			}
		}

		if paused {
			continue
		}

		for i := 0; i < c.opts.TPS; i++ {
			if err := c.vm.RunCycle(); err != nil {
				log.Error(err)
			}
		}

		c.vm.UpdateTimers()

		c.updateScreen(renderer, &rect)
	}

	log.Info("Quitting CHIP-8")

	return nil
}

func (c *Chip8) handleEvent() Event {
	if event := sdl.PollEvent(); event != nil {
		switch event.GetType() {
		case sdl.QUIT:
			return EventQuit

		case sdl.KEYUP, sdl.KEYDOWN:
			return c.handleKeyEvent(event.(*sdl.KeyboardEvent))
		}
	}

	return EventIgnore
}

func (c *Chip8) handleKeyEvent(event *sdl.KeyboardEvent) Event {
	pressed := event.GetType() == sdl.KEYDOWN

	switch scancode := event.Keysym.Scancode; scancode {
	case sdl.SCANCODE_SPACE:
		if pressed {
			return EventPause
		}

	case sdl.SCANCODE_ESCAPE:
		if pressed {
			return EventQuit
		}

	case sdl.SCANCODE_SEMICOLON:
		if pressed {
			return EventReset
		}

	case sdl.SCANCODE_COMMA:
		if pressed {
			return EventNextCycle
		}

	default:
		if key, ok := scancodeHexMap[scancode]; ok {
			c.keyBuffer[key] = pressed

			if pressed && c.lastKeyPress == -1 {
				c.lastKeyPress = int(key)
			}
		}
	}

	return EventIgnore
}

func (c Chip8) updateScreen(renderer *sdl.Renderer, rect *sdl.Rect) {
	renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	renderer.Clear()
	renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	for i, on := range c.screenBuffer {
		if on {
			rect.X = int32(i%ScreenWidth) * rect.W
			rect.Y = int32(i/ScreenWidth) * rect.H
			renderer.FillRect(rect)
		}
	}

	renderer.Present()
}

//
// Keyboard
//

func (c *Chip8) IsKeyPressed(key uint8) bool {
	return c.keyBuffer[key]
}

func (c *Chip8) PollKeyPress() (uint8, bool) {
	if c.lastKeyPress != -1 {
		key := uint8(c.lastKeyPress)
		c.lastKeyPress = -1
		return key, true
	}
	return 0, false
}

//
// Screen
//

func (c *Chip8) ClearScreen() {
	for i := 0; i < len(c.screenBuffer); i++ {
		c.screenBuffer[i] = false
	}
}

func (c *Chip8) SetSprite(x, y uint8, sprite ...uint8) bool {
	flag := false

	for i, b := range sprite {
		for j := uint8(0); j < 8; j++ {
			idx := getScreenBufferIndex(x+(7-j), y+uint8(i))
			bit := b&1 == 1

			if c.screenBuffer[idx] && bit {
				flag = true
			}

			c.screenBuffer[idx] = c.screenBuffer[idx] != bit
			b >>= 1
		}
	}

	return flag
}

func getScreenBufferIndex(x, y uint8) int {
	return (int(y%ScreenHeight) * ScreenWidth) + int(x%ScreenWidth)
}
