package ch8

import "time"

const (
	// DefaultHertzIO is the default speed (in hertz) in which to update
	// the IO timers and audio.
	DefaultHertzIO = 16 * time.Millisecond

	// DefaultHertzVM is the default speed (in hertz) in which to run a
	// CPU cycle of the CHIP-8 virtual machine.
	DefaultHertzVM = 2 * time.Millisecond

	// DefaultMaxTPS is the default max ticks-per-second (TPS) of the
	// renderer.
	DefaultMaxTPS = 60

	// DefaultScale is the default scale factor of the CHIP-8 screen.
	DefaultScale = 10
)
