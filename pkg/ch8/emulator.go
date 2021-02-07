package ch8

import (
	"image/color"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	vm *VirtualMachine
}

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator() *Emulator {
	return &Emulator{
		vm: NewVirtualMachine(),
	}
}

// Start starts the emulator.
func (emu *Emulator) Start() error {
	ebiten.SetWindowSize(DisplayWidth*10, DisplayHeight*10)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetMaxTPS(60)
	ebiten.SetVsyncEnabled(true)

	go emu.startVM()
	go emu.startTimers()

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
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			if emu.vm.Display[y*DisplayWidth+x] == 0x1 {
				screen.Set(x, y, fg)
			} else {
				screen.Set(x, y, bg)
			}
		}
	}
}

// Layout returns the resolution of the emulator's screen.
func (emu *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func (emu *Emulator) startVM() {
	// Run VM separately at approximately 500 Hz
	for range time.Tick(2 * time.Millisecond) {
		if err := emu.vm.RunCycle(); err != nil {
			log.Println(err)
		}
	}
}

func (emu *Emulator) startTimers() {
	// Run timers separately at approximately 60 Hz
	for range time.Tick(16 * time.Millisecond) {
		emu.vm.UpdateTimers()
	}
}
