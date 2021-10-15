package ch8

import "time"

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
