package ch8

import (
	"errors"
	"fmt"
)

//===========================================================================
// Emulator
//===========================================================================

var ErrTerminated = errors.New("Emulator terminated")

//===========================================================================
// Virtual machine
//===========================================================================

// InvalidProgramError is an error that occurs from loading a program.
func InvalidProgramError(msg string) error {
	return fmt.Errorf("invalid program: %s", msg)
}

// InvalidStateError is an error that occurs due to bad state in the
// virtual machine.
func InvalidStateError(msg string) error {
	return fmt.Errorf("invalid state: %s", msg)
}

// InvalidJumpError is an error caused by a program trying to jump to
// an invalid memory location.
func InvalidJumpError(fromAddr, toAddr uint16) error {
	return fmt.Errorf("invalid jump: Jump from %.3X to %.3X", fromAddr, toAddr)
}
