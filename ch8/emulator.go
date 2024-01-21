package ch8

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// EmulatorOptions is a set of arguments that allows you to set
// different options in the emulator.
type EmulatorOptions struct {
	// Scale is the scale factor of the CHIP-8 screen.
	Scale int
}

const (
	eventPlay      EmulatorEvent = "play"
	eventPause     EmulatorEvent = "pause"
	eventReset     EmulatorEvent = "reset"
	eventTerminate EmulatorEvent = "terminate"
)

var (
	foreground = color.White
	background = color.Black

	keyHexMap = map[ebiten.Key]uint{
		ebiten.Key1: 0x0, ebiten.Key2: 0x1, ebiten.Key3: 0x2, ebiten.Key4: 0x3,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0x7,
		ebiten.KeyA: 0x8, ebiten.KeyS: 0x9, ebiten.KeyD: 0xa, ebiten.KeyF: 0xb,
		ebiten.KeyZ: 0xc, ebiten.KeyX: 0xd, ebiten.KeyC: 0xe, ebiten.KeyV: 0xf,
	}

	keyEventMap = map[ebiten.Key]EmulatorEvent{
		ebiten.KeyRightBracket: eventPause,
		ebiten.KeyLeftBracket:  eventPlay,
		ebiten.KeyBackslash:    eventReset,
		ebiten.KeyEscape:       eventTerminate,
	}
)

// Emulator is the CHIP-8 emulator.
type Emulator struct {
	scale     int
	vm        *VirtualMachine
	vmChannel chan EmulatorEvent
}

// EmulatorEvent is an event that occurs that controls the state of the
// emulator.
type EmulatorEvent string

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator(opts EmulatorOptions) (*Emulator, error) {
	if opts.Scale < 1 {
		return nil, errors.New("scale factor must be positive")
	}

	return &Emulator{
		scale:     opts.Scale,
		vm:        NewVirtualMachine(),
		vmChannel: make(chan EmulatorEvent),
	}, nil
}

// Start starts the emulator.
func (emu *Emulator) Start() (err error) {
	ebiten.SetWindowSize(DisplayWidth*emu.scale, DisplayHeight*emu.scale)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetTPS(DefaultMaxTPS)
	ebiten.SetVsyncEnabled(true)

	go emu.startIO()
	go emu.startVM()

	return ebiten.RunGame(emu)
}

// LoadROM reads a CHIP-8 ROM program file (*.ch8) and loads it into
// the virtual machine's memory.
func (emu *Emulator) LoadROM(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return emu.LoadBytes(data)
}

// Update updates the state of the emulator.
func (emu *Emulator) Update() error {
	for key, event := range keyEventMap {
		if ebiten.IsKeyPressed(key) {
			if event == eventTerminate {
				return ErrTerminated
			}

			emu.vmChannel <- event
			return nil
		}
	}

	for key, hex := range keyHexMap {
		emu.vm.Keys[hex] = ebiten.IsKeyPressed(key)
	}
	return nil
}

// Draw renders the screen of the emulator.
func (emu *Emulator) Draw(screen *ebiten.Image) {
	screen.Fill(background)

	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			if emu.vm.Display[y][x] {
				screen.Set(x, y, foreground)
			}
		}
	}

	ebiten.SetWindowTitle(fmt.Sprintf("CHIP-8 | FPS: %.2f", ebiten.ActualFPS()))
}

// Layout returns the resolution of the emulator's screen.
func (emu *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func (emu *Emulator) startVM() {
	pause := false

	for range time.Tick(DefaultHertzVM) {
		select {
		case event := <-emu.vmChannel:
			switch event {
			case eventPlay:
				pause = false
			case eventPause:
				pause = true
			case eventReset:
				emu.vm.Reset()
			}
		default:
			if pause {
				continue
			}

			if err := emu.vm.RunCycle(); err != nil {
				log.Println(err)
			}
		}
	}
}

func (emu *Emulator) startIO() {
	for range time.Tick(DefaultHertzIO) {
		emu.vm.UpdateTimers()
	}
}
