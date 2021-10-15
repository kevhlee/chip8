package ch8

import (
	"errors"
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

	// DefaultHertzIO is the default speed (in hertz) in which to update
	// the IO timers and audio.
	DefaultHertzIO = 16 * time.Millisecond

	// DefaultHertzVM is the default speed (in hertz) in which to run a
	// CPU cycle of the CHIP-8 virtual machine.
	DefaultHertzVM = 2 * time.Millisecond

	// DefaultMaxTPS is the default max ticks-per-second (TPS) of the
	// renderer.
	DefaultMaxTPS = 60

	// DefaultSampleRate is the default sample rate of the CHIP-8
	// beeper.
	DefaultSampleRate = 44100

	// DefaultScale is the default scale factor of the CHIP-8 screen.
	DefaultScale = 10

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
	EventPlay  EmulatorEvent = "play"
	EventPause EmulatorEvent = "pause"
	EventReset EmulatorEvent = "reset"
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
		ebiten.KeyRightBracket: EventPause,
		ebiten.KeyLeftBracket:  EventPlay,
		ebiten.KeyBackslash:    EventReset,
	}
)

// Emulator is the CHIP-8 emulator.
type Emulator struct {
	beeper    *audio.Player
	options   *EmulatorOptions
	vm        *VirtualMachine
	vmChannel chan EmulatorEvent
}

// EmulatorEvent is an event that occurs that controls the state of the
// emulator.
type EmulatorEvent string

// EmulatorOptions is a set of arguments that allows you to set
// different options in the emulator.
type EmulatorOptions struct {
	// Frequency is the frequency of the CHIP-8 beeper.
	Frequency int

	// HertzIO is the speed (in hertz) in which to update the IO timers
	// and audio.
	HertzIO time.Duration

	// HertzVM is the speed in which to run a CPU cycle of the CHIP-8
	// virtual machine.
	HertzVM time.Duration

	// MaxTPS is the max ticks-per-second (TPS) of the renderer.
	MaxTPS int

	// SampleRate is the sample rate of the CHIP-8 beeper.
	SampleRate int

	// Scale is the scale factor of the CHIP-8 screen.
	Scale int

	// Volume is the volume of the CHIP-8 beeper.
	//
	// The volume ranges within [0.0, 1.0].
	Volume float64
}

// NewEmulator creates a new CHIP-8 emulator instance.
func NewEmulator(options *EmulatorOptions) (*Emulator, error) {
	if err := checkOptions(options); err != nil {
		return nil, err
	}

	return &Emulator{
		options:   options,
		vm:        NewVirtualMachine(),
		vmChannel: make(chan EmulatorEvent),
	}, nil
}

func NewEmulatorOptions() *EmulatorOptions {
	return &EmulatorOptions{
		Frequency:  DefaultFrequency,
		HertzIO:    DefaultHertzIO,
		HertzVM:    DefaultHertzVM,
		MaxTPS:     DefaultMaxTPS,
		SampleRate: DefaultSampleRate,
		Scale:      DefaultScale,
		Volume:     DefaultVolume,
	}
}

// Start starts the emulator.
func (emu *Emulator) Start() (err error) {
	ebiten.SetWindowSize(DisplayWidth*emu.options.Scale, DisplayHeight*emu.options.Scale)
	ebiten.SetWindowTitle("CHIP-8")
	ebiten.SetMaxTPS(emu.options.MaxTPS)
	ebiten.SetVsyncEnabled(true)

	emu.beeper, err = audio.NewPlayer(
		audio.NewContext(emu.options.SampleRate),
		&stream{
			frequency:  emu.options.Frequency,
			sampleRate: emu.options.SampleRate,
		},
	)
	if err != nil {
		return err
	}

	emu.beeper.SetVolume(emu.options.Volume)

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

	ebiten.SetWindowTitle(
		fmt.Sprintf("CHIP-8 | FPS: %.2f", ebiten.CurrentFPS()),
	)
}

// Layout returns the resolution of the emulator's screen.
func (emu *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return DisplayWidth, DisplayHeight
}

func checkOptions(options *EmulatorOptions) error {
	if options.Scale < 1 {
		return errors.New("scale factor must be positive")
	}

	if options.Volume < 0.0 || options.Volume > 1.0 {
		return errors.New("volume must be between [0, 1]")
	}

	return nil
}

func (emu *Emulator) startVM() {
	pause := false

	for range time.Tick(emu.options.HertzVM) {
		select {
		case event := <-emu.vmChannel:
			switch event {
			case EventPlay:
				pause = false
			case EventPause:
				pause = true
			case EventReset:
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
	for range time.Tick(emu.options.HertzIO) {
		emu.vm.UpdateTimers()

		if emu.vm.ST > 0x00 {
			emu.beeper.Play()
		} else {
			emu.beeper.Pause()
		}
	}
}
