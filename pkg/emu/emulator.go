package emu

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kevhlee/chip8/pkg/ch8"
)

var (
	fg = color.RGBA{0x00, 0xff, 0x00, 0xff}
	bg = color.RGBA{0x00, 0x00, 0x00, 0xff}

	keymap = map[ebiten.Key]uint{
		ebiten.Key1: 0x0, ebiten.Key2: 0x1, ebiten.Key3: 0x2, ebiten.Key4: 0x3,
		ebiten.KeyQ: 0x4, ebiten.KeyW: 0x5, ebiten.KeyE: 0x6, ebiten.KeyR: 0x7,
		ebiten.KeyA: 0x8, ebiten.KeyS: 0x9, ebiten.KeyD: 0xa, ebiten.KeyF: 0xb,
		ebiten.KeyZ: 0xc, ebiten.KeyX: 0xd, ebiten.KeyC: 0xe, ebiten.KeyV: 0xf,
	}
)

// Emulator is the CHIP-8 emulator.
type Emulator struct {
	vm     *ch8.VirtualMachine
	beeper *Beeper
	mute   bool
}

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator(scale int, mute bool) *Emulator {
	ebiten.SetWindowSize(ch8.DisplayWidth*scale, ch8.DisplayHeight*scale)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetMaxTPS(60)
	ebiten.SetVsyncEnabled(true)

	return &Emulator{
		vm:     ch8.NewVirtualMachine(),
		beeper: NewBeeper(DefaultFrequency, DefaultSampleRate),
		mute:   mute,
	}
}

// Start starts the emulator.
func (emu *Emulator) Start() error {
	go emu.startVM()
	go emu.startBeeper()
	return ebiten.RunGame(emu)
}

// LoadROM loads a CHIP-8 ROM into the virtual machine.
func (emu *Emulator) LoadROM(path string) error {
	return emu.vm.LoadROM(path)
}

// Update updates the state of the emulator.
func (emu *Emulator) Update() error {
	emu.handleKeys()
	emu.vm.UpdateTimers()
	return nil
}

// Draw renders the screen of the emulator.
func (emu *Emulator) Draw(screen *ebiten.Image) {
	for y := 0; y < ch8.DisplayHeight; y++ {
		for x := 0; x < ch8.DisplayWidth; x++ {
			if emu.vm.Display[y][x] {
				screen.Set(x, y, fg)
			} else {
				screen.Set(x, y, bg)
			}
		}
	}
}

// Layout returns the resolution of the emulator's screen.
func (emu *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ch8.DisplayWidth, ch8.DisplayHeight
}

func (emu *Emulator) startVM() {
	for range time.Tick(2 * time.Millisecond) {
		emu.vm.RunCycle()
	}
}

func (emu *Emulator) startBeeper() {
	if emu.mute {
		return
	}

	for range time.Tick(2 * time.Millisecond) {
		if emu.vm.ST > 0x00 {
			emu.beeper.Play()
		} else {
			emu.beeper.Stop()
		}
	}
}

func (emu *Emulator) handleKeys() {
	for key, hex := range keymap {
		emu.vm.Keys[hex] = ebiten.IsKeyPressed(key)
	}
}
