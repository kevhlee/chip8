package ch8

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
)

//=====================================================================
// Constants
//=====================================================================

const (
	// DefaultFrequency is the default frequency of the CHIP-8 beeper.
	DefaultFrequency = 440

	// DefaultHertzI is the default speed (in hertz) in which to update
	// the IO timers and audio.
	DefaultHertzIO = 16 * time.Millisecond

	// DefaultHertzVM is the default speed (in hertz) in which to run a
	// CPU cycle of the CHIP-8 virtual machine.
	DefaultHertzVM = 2 * time.Millisecond

	// DefaultSampleRate is the default sample rate of the CHIP-8
	// beeper.
	DefaultSampleRate = 44100

	// DefaultScale is the default scale factor of the CHIP-8 screen.
	DefaultScale = 10

	// DefaultTPS is the default ticks per second of the emulator.
	DefaultTPS = 60

	// DefaultVolume is the default volume of the CHIP-8 beeper.
	//
	// The volume ranges within [0.0, 1.0].
	DefaultVolume = 0.5
)

//=====================================================================
// Audio
//=====================================================================

// This struct is taken directly from Ebiten's example code:
// <https://ebiten.org/examples/sinewave.html>
type stream struct {
	frequency  int
	sampleRate int
	position   int64
	remaining  []byte
}

// Read fills the byte stream with sine wave samples.
func (s *stream) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	length := int64(s.sampleRate / s.frequency)
	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * max)
		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)
		p++
	}

	s.position += int64(len(buf))
	s.position %= length * 4

	if origBuf != nil {
		n := copy(origBuf, buf)
		s.remaining = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

// Close closes the bye stream.
func (s *stream) Close() error {
	return nil
}

//=====================================================================
// Emulator
//=====================================================================

const (
	playEvent  = "play"
	pauseEvent = "pause"
	resetEvent = "reset"
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
	keyEventMap = map[ebiten.Key]string{
		ebiten.KeyRightBracket: pauseEvent,
		ebiten.KeyLeftBracket:  playEvent,
		ebiten.KeyBackslash:    resetEvent,
	}
)

// Emulator is the CHIP-8 emulator.
type Emulator struct {
	vm     *VirtualMachine
	beeper *audio.Player
	vmChan chan string
}

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator(scale int, volume float64) *Emulator {
	// Initialize audio
	beeper, _ := audio.NewPlayer(
		audio.NewContext(DefaultSampleRate),
		&stream{
			frequency:  DefaultFrequency,
			sampleRate: DefaultSampleRate,
		},
	)
	beeper.SetVolume(volume)

	// Initialize graphics
	ebiten.SetWindowSize(DisplayWidth*scale, DisplayHeight*scale)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetMaxTPS(DefaultTPS)
	ebiten.SetVsyncEnabled(true)

	return &Emulator{NewVirtualMachine(), beeper, make(chan string)}
}

// Start starts the emulator.
func (emu *Emulator) Start() error {
	go emu.startIO()
	go emu.startVM()

	return ebiten.RunGame(emu)
}

// LoadROM loads a CHIP-8 ROM into the virtual machine.
func (emu *Emulator) LoadROM(path string) error {
	return emu.vm.LoadROM(path)
}

// Update updates the state of the emulator.
func (emu *Emulator) Update() error {
	for key, event := range keyEventMap {
		if ebiten.IsKeyPressed(key) {
			emu.vmChan <- event
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

	ebiten.SetWindowTitle(
		fmt.Sprintf("CHIP-8 | FPS: %.2f", ebiten.CurrentFPS()),
	)
}

// Layout returns the resolution of the emulator's screen.
func (emu *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func (emu *Emulator) startVM() {
	pause := false

	for range time.Tick(DefaultHzVM) {
		select {
		case event := <-emu.vmChan:
			switch event {
			case playEvent:
				pause = false
			case pauseEvent:
				pause = true
			case resetEvent:
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
	for range time.Tick(DefaultHzIO) {
		emu.vm.UpdateTimers()

		if emu.vm.ST > 0x00 {
			emu.beeper.Play()
		} else {
			emu.beeper.Pause()
		}
	}
}
