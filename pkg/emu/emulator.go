package emu

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kevhlee/chip8/pkg/ch8"
)

var (
	foreground = color.RGBA{0x00, 0xff, 0x00, 0xff}
	background = color.RGBA{0x00, 0x00, 0x00, 0xff}
	keymap     = map[ebiten.Key]uint{
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
	debug  bool
	scale  int
}

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator(debug bool, scale int, mute bool) *Emulator {
	return &Emulator{
		vm:     ch8.NewVirtualMachine(),
		beeper: NewBeeper(DefaultFrequency, DefaultSampleRate),
		mute:   mute,
		debug:  debug,
		scale:  scale,
	}
}

// Start starts the emulator.
func (emu *Emulator) Start() error {
	ebiten.SetWindowSize(
		ch8.DisplayWidth*emu.scale,
		ch8.DisplayHeight*emu.scale,
	)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetMaxTPS(60)
	ebiten.SetVsyncEnabled(true)

	go emu.startTimers()
	go emu.startBeeper()
	go emu.startVM()

	return ebiten.RunGame(emu)
}

// LoadROM loads a CHIP-8 ROM into the virtual machine.
func (emu *Emulator) LoadROM(path string) error {
	return emu.vm.LoadROM(path)
}

// Update updates the state of the emulator.
func (emu *Emulator) Update() error {
	for key, hex := range keymap {
		emu.vm.Keys[hex] = ebiten.IsKeyPressed(key)
	}
	return nil
}

// Draw renders the screen of the emulator.
func (emu *Emulator) Draw(screen *ebiten.Image) {
	for y := 0; y < ch8.DisplayHeight; y++ {
		for x := 0; x < ch8.DisplayWidth; x++ {
			if emu.vm.Display[y][x] {
				screen.Set(x, y, foreground)
			} else {
				screen.Set(x, y, background)
			}
		}
	}

	ebiten.SetWindowTitle(
		fmt.Sprintf("CHIP-8 | FPS: %.2f", ebiten.CurrentFPS()),
	)
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

func (emu *Emulator) startTimers() {
	for range time.Tick(16 * time.Millisecond) {
		emu.vm.UpdateTimers()
	}
}

func (emu *Emulator) startBeeper() {
	if emu.mute {
		return
	}

	for range time.Tick(16 * time.Millisecond) {
		if emu.vm.ST > 0x00 {
			emu.beeper.Play()
		} else {
			emu.beeper.Stop()
		}
	}
}
