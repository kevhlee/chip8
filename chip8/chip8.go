package chip8

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/veandco/go-sdl2/sdl"
)

type Options struct {
	Scale int
	TPS   int
}

type Event int

const (
	EventIgnore Event = iota
	EventQuit
	EventPause
	EventReset
	EventNextCycle
)

var (
	scancodeMap map[sdl.Scancode]uint8 = map[sdl.Scancode]uint8{
		sdl.SCANCODE_1: 0x1, sdl.SCANCODE_2: 0x2, sdl.SCANCODE_3: 0x3, sdl.SCANCODE_4: 0xC,
		sdl.SCANCODE_Q: 0x4, sdl.SCANCODE_W: 0x5, sdl.SCANCODE_E: 0x6, sdl.SCANCODE_R: 0xD,
		sdl.SCANCODE_A: 0x7, sdl.SCANCODE_S: 0x8, sdl.SCANCODE_D: 0x9, sdl.SCANCODE_F: 0xE,
		sdl.SCANCODE_Z: 0xA, sdl.SCANCODE_X: 0x0, sdl.SCANCODE_C: 0xB, sdl.SCANCODE_V: 0xF,
	}
)

func Start(filename string, opts Options) error {
	log.Info("Initializing CHIP-8")

	if err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO); err != nil {
		return err
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"CHIP-8",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(ScreenWidth*opts.Scale),
		int32(ScreenHeight*opts.Scale),
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

	var (
		paused  = false
		running = true

		vm       = NewVirtualMachine()
		keyboard = NewKeyboard()
		screen   = NewScreen(renderer, int32(opts.Scale))
		sound    = NewSound()
		timer    = NewTimer()
	)

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if err := vm.LoadProgram(data...); err != nil {
		return err
	}

	log.Info("Starting CHIP-8")

	// Runs at roughly 60 FPS
	for range time.Tick(time.Millisecond * 1000 / 60) {
		if !running {
			break
		}

		switch handleEvent(keyboard) {
		case EventQuit:
			running = false

		case EventPause:
			paused = !paused

		case EventReset:
			vm.Reset()
			screen.Render()

		case EventNextCycle:
			if paused {
				vm.Step(keyboard, screen, sound, timer)
				screen.Render()
			}
		}

		if paused {
			continue
		}

		sound.Step()
		timer.Step()

		for i := 0; i < opts.TPS; i++ {
			if err := vm.Step(keyboard, screen, sound, timer); err != nil {
				log.Error(err)
			}
		}

		screen.Render()
	}

	return nil
}

func handleEvent(keyboard *Keyboard) Event {
	if event := sdl.PollEvent(); event != nil {
		switch event.GetType() {
		case sdl.QUIT:
			return EventQuit

		case sdl.KEYUP, sdl.KEYDOWN:
			return handleKeyEvent(event.(*sdl.KeyboardEvent), keyboard)
		}
	}

	return EventIgnore
}

func handleKeyEvent(event *sdl.KeyboardEvent, keyboard *Keyboard) Event {
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
		if key, ok := scancodeMap[scancode]; ok {
			keyboard.Set(key, pressed)
		}
	}

	return EventIgnore
}
