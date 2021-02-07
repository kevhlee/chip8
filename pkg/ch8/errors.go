package ch8

import (
	"fmt"
)

// InvalidProgramError is an error that occurs from loading a program.
func InvalidProgramError(msg string) error {
	return fmt.Errorf("Invalid program error: %s", msg)
}

// InvalidStateError is an error that occurs when an issue occurs in
// the virtual machine.
func InvalidStateError(msg string) error {
	return fmt.Errorf("Invalid state error: %s", msg)
}

// InvalidJumpError is an error caused by a program trying to jump to
// an invalid memory location.
func InvalidJumpError(fromAddr, toAddr uint) error {
	return fmt.Errorf("Invalid jump error: Jump from %.3X to %.3X", fromAddr, toAddr)
}

// InvalidOpcodeError is an error caused by the CHIP-8 virtual machine
// trying to run an invalid opcode.
func InvalidOpcodeError(opcode uint) error {
	return fmt.Errorf("Invalid opcode error: %.4X", opcode)
}
