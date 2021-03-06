package ch8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
)

//===========================================================================
// Constants
//===========================================================================

const (
	// FontSize is the number of bytes in a CHIP-8 built-in font.
	FontSize = 0x5

	// MaxStackDepth is the maximum depth of the CHIP-8 virtual
	// machine's call stack.
	MaxStackDepth = 0x10

	// MemorySize is the total amount of memory available in the CHIP-8
	// virtual machine.
	MemorySize = 0x1000

	// DisplayWidth is the width (in pixels) of the CHIP-8 display.
	DisplayWidth = 0x40

	// DisplayHeight is the height (in pixels) of the CHIP-8 display.
	DisplayHeight = 0x20

	// NumberOfKeys is the number of keys in the CHIP-8 keyboard.
	NumberOfKeys = 0x10

	// NumberOfFonts is the total number of built-in fonts in the
	// CHIP-8 virtual machine.
	NumberOfFonts = 0x10

	// NumberOfPixels is the total number of pixels in the CHIP-8
	// display.
	NumberOfPixels = DisplayWidth * DisplayHeight

	// NumberOfRegisters is the number of general-purpose registers in
	// the CHIP-8 virtual machine.
	NumberOfRegisters = 0x10

	// ProgramStartAddress is the start memory location where programs
	// are loaded in the virtual machine's memory.
	ProgramStartAddress = 0x200

	// ProgramMemorySize is the total amount of memory available for
	// CHIP-8 programs.
	ProgramMemorySize = MemorySize - ProgramStartAddress
)

//===========================================================================
// Errors
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
func InvalidJumpError(fromAddr, toAddr uint) error {
	return fmt.Errorf(
		"invalid jump: Jump from %.3X to %.3X",
		fromAddr,
		toAddr,
	)
}

// InvalidOpcodeError is an error caused by the CHIP-8 virtual machine
// trying to run an invalid opcode.
func InvalidOpcodeError(opcode uint) error {
	return fmt.Errorf("invalid opcode: %.4X", opcode)
}

//===========================================================================
// Virtual Machine
//===========================================================================

// VirtualMachine is the CHIP-8 virtual machine.
type VirtualMachine struct {
	I        uint
	SP       uint
	PC       uint
	DT       uint
	ST       uint
	V        [NumberOfRegisters]uint
	Stack    [MaxStackDepth]uint
	Memory   [MemorySize]uint
	Keys     [NumberOfKeys]bool
	Display  [DisplayHeight][DisplayWidth]bool
	Opcode   uint
	opcodeFn map[uint]func() error
}

// NewVirtualMachine creates new CHIP-8 virtual machine instance.
func NewVirtualMachine() *VirtualMachine {
	vm := &VirtualMachine{
		PC:      ProgramStartAddress,
		Stack:   [MaxStackDepth]uint{},
		V:       [NumberOfRegisters]uint{},
		Keys:    [NumberOfKeys]bool{},
		Display: [DisplayHeight][DisplayWidth]bool{},
		Memory:  [MemorySize]uint{},
	}

	fonts := []uint{
		0xf0, 0x90, 0x90, 0x90, 0xf0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xf0, 0x10, 0xf0, 0x80, 0xf0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xf0, // 3
		0x90, 0x90, 0xf0, 0x10, 0x10, // 4
		0xf0, 0x80, 0xf0, 0x10, 0xf0, // 5
		0xf0, 0x80, 0xf0, 0x90, 0xf0, // 6
		0xf0, 0x10, 0x20, 0x40, 0x40, // 7
		0xf0, 0x90, 0xf0, 0x90, 0xf0, // 8
		0xf0, 0x90, 0xf0, 0x10, 0xf0, // 9
		0xf0, 0x90, 0xf0, 0x90, 0x90, // A
		0xe0, 0x90, 0xe0, 0x90, 0xe0, // B
		0xf0, 0x80, 0x80, 0x80, 0xf0, // C
		0xe0, 0x90, 0x90, 0x90, 0xe0, // D
		0xf0, 0x80, 0xf0, 0x80, 0xf0, // E
		0xf0, 0x80, 0xf0, 0x80, 0x80, // F
	}

	for i, b := range fonts {
		vm.Memory[i] = b
	}

	vm.opcodeFn = map[uint]func() error{
		0x0: vm.executeOp0x0, 0x1: vm.executeOp0x1,
		0x2: vm.executeOp0x2, 0x3: vm.executeOp0x3,
		0x4: vm.executeOp0x4, 0x5: vm.executeOp0x5,
		0x6: vm.executeOp0x6, 0x7: vm.executeOp0x7,
		0x8: vm.executeOp0x8, 0x9: vm.executeOp0x9,
		0xa: vm.executeOp0xA, 0xb: vm.executeOp0xB,
		0xc: vm.executeOp0xC, 0xd: vm.executeOp0xD,
		0xe: vm.executeOp0xE, 0xf: vm.executeOp0xF,
	}

	return vm
}

// RunCycle runs a single CPU cycle of the virtual machine.
func (vm *VirtualMachine) RunCycle() error {
	// Fetch-decode-execute
	vm.fetch()
	execute := vm.decode()
	err := execute()

	// Keep program counter within range
	if vm.PC > 0xfff {
		vm.PC = (vm.PC & 0xfff) + ProgramStartAddress
	}

	return err
}

// UpdateTimers updates the delay and sound timers.
func (vm *VirtualMachine) UpdateTimers() {
	if vm.DT > 0x00 {
		vm.DT--
	}
	if vm.ST > 0x00 {
		vm.ST--
	}
}

// LoadROM reads a CHIP-8 ROM program file (*.ch8) and loads it into
// memory.
func (vm *VirtualMachine) LoadROM(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	} else if len(data) > ProgramMemorySize {
		return InvalidProgramError("The ROM is too large")
	}

	i := ProgramStartAddress
	for _, b := range data {
		vm.Memory[i] = uint(b)
		i++
	}

	return nil
}

// LoadOpcodes loads opcodes into the virtual machine's program memory.
func (vm *VirtualMachine) LoadOpcodes(opcodes []uint) error {
	if 2*len(opcodes) >= ProgramMemorySize {
		return InvalidProgramError("The ROM is too large")
	}

	i := ProgramStartAddress
	for _, opcode := range opcodes {
		if opcode > 0xffff {
			return InvalidOpcodeError(vm.Opcode)
		}

		vm.Memory[i] = opcode >> 8
		vm.Memory[i+1] = opcode & 0xff
		i += 0x2
	}

	return nil
}

// Reset resets the virtual machine.
//
// This preserves the program/opcodes already loaded in memory.
func (vm *VirtualMachine) Reset() {
	vm.ClearRegisters()
	vm.ClearDisplay()
	vm.ClearKeys()
}

// Clear clears the entire state of the virtual machine.
func (vm *VirtualMachine) Clear() {
	vm.Reset()
	vm.ClearProgram()
}

// ClearKeys clears the state of the keys.
func (vm *VirtualMachine) ClearKeys() {
	for i := 0; i < len(vm.Keys); i++ {
		vm.Keys[i] = false
	}
}

// ClearProgram clears the program loaded in the virtual machine.
func (vm *VirtualMachine) ClearProgram() {
	for i := ProgramStartAddress; i < len(vm.Memory); i++ {
		vm.Memory[i] = 0x00
	}
}

// ClearRegisters clears all the registers, including the program
// counter, timers, stack pointer.
func (vm *VirtualMachine) ClearRegisters() {
	vm.I = 0x000
	vm.SP = 0x00
	vm.PC = ProgramStartAddress
	vm.DT = 0x00
	vm.ST = 0x00

	for i := 0; i < len(vm.V); i++ {
		vm.V[i] = 0x00
	}
}

// ClearDisplay clears the state of the display.
func (vm *VirtualMachine) ClearDisplay() {
	for y := 0; y < DisplayHeight; y++ {
		for x := 0; x < DisplayWidth; x++ {
			vm.Display[y][x] = false
		}
	}
}

//=====================================================================
// CPU Cycle
//=====================================================================

func (vm *VirtualMachine) fetch() {
	opcode := (vm.Memory[vm.PC] << 8) | vm.Memory[vm.PC+1]
	vm.PC += 0x2
	vm.Opcode = opcode
}

func (vm *VirtualMachine) decode() func() error {
	return vm.opcodeFn[vm.decodeOp()]
}

func (vm *VirtualMachine) decodeX() uint {
	return (vm.Opcode >> 8) & 0xf
}

func (vm *VirtualMachine) decodeY() uint {
	return (vm.Opcode >> 4) & 0xf
}

func (vm *VirtualMachine) decodeN() uint {
	return vm.Opcode & 0xf
}

func (vm *VirtualMachine) decodeOp() uint {
	return vm.Opcode >> 0xc
}

func (vm *VirtualMachine) decodeKK() uint {
	return vm.Opcode & 0xff
}

func (vm *VirtualMachine) decodeNNN() uint {
	return vm.Opcode & 0xfff
}

func (vm *VirtualMachine) executeOp0x0() error {
	switch vm.decodeNNN() {
	case 0x0e0:
		vm.ClearDisplay()
	case 0x0ee:
		vm.SP--
		vm.PC = vm.Stack[vm.SP]
	default:
		return InvalidOpcodeError(vm.Opcode)
	}
	return nil
}

func (vm *VirtualMachine) executeOp0x1() error {
	nnn := vm.decodeNNN()

	if nnn < ProgramStartAddress {
		return InvalidJumpError(vm.PC, nnn)
	}

	vm.PC = nnn

	return nil
}

func (vm *VirtualMachine) executeOp0x2() error {
	nnn := vm.decodeNNN()

	if vm.SP >= MaxStackDepth {
		return InvalidStateError("Stack overflow")
	} else if nnn < ProgramStartAddress {
		return InvalidJumpError(vm.PC, nnn)
	}

	vm.Stack[vm.SP] = vm.PC
	vm.PC = nnn
	vm.SP++

	return nil
}

func (vm *VirtualMachine) executeOp0x3() error {
	if vm.V[vm.decodeX()] == vm.decodeKK() {
		vm.PC += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp0x4() error {
	if vm.V[vm.decodeX()] != vm.decodeKK() {
		vm.PC += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp0x5() error {
	if vm.decodeN() != 0x0 {
		return InvalidOpcodeError(vm.Opcode)
	}

	if vm.V[vm.decodeX()] == vm.V[vm.decodeY()] {
		vm.PC += 2
	}

	return nil
}

func (vm *VirtualMachine) executeOp0x6() error {
	vm.V[vm.decodeX()] = vm.decodeKK()
	return nil
}

func (vm *VirtualMachine) executeOp0x7() error {
	x := vm.decodeX()
	kk := vm.decodeKK()
	vm.V[x] = (vm.V[x] + kk) & 0xff
	return nil
}

func (vm *VirtualMachine) executeOp0x8() error {
	x := vm.decodeX()
	y := vm.decodeY()

	switch vm.decodeN() {
	case 0x0:
		vm.V[x] = vm.V[y]
	case 0x1:
		vm.V[x] |= vm.V[y]
	case 0x2:
		vm.V[x] &= vm.V[y]
	case 0x3:
		vm.V[x] ^= vm.V[y]
	case 0x4:
		result := vm.V[x] + vm.V[y]
		if result > 0xff {
			vm.V[0xf] = 0x1
			vm.V[x] = result & 0xff
		} else {
			vm.V[0xf] = 0x0
			vm.V[x] = result
		}
	case 0x5:
		if vm.V[x] > vm.V[y] {
			vm.V[0xf] = 0x1
		} else {
			vm.V[0xf] = 0x0
		}
		vm.V[x] = (vm.V[x] - vm.V[y]) & 0xff
	case 0x6:
		vm.V[0xf] = vm.V[x] & 0x01
		vm.V[x] >>= 1
	case 0x7:
		if vm.V[x] < vm.V[y] {
			vm.V[0xf] = 0x1
		} else {
			vm.V[0xf] = 0x0
		}
		vm.V[x] = (vm.V[y] - vm.V[x]) & 0xff
	case 0xe:
		vm.V[0xf] = vm.V[x] >> 7
		vm.V[x] = (vm.V[x] << 1) & 0xff
	default:
		return InvalidOpcodeError(vm.Opcode)
	}

	return nil
}

func (vm *VirtualMachine) executeOp0x9() error {
	if vm.decodeN() != 0x0 {
		return InvalidOpcodeError(vm.Opcode)
	}

	if vm.V[vm.decodeX()] != vm.V[vm.decodeY()] {
		vm.PC += 2
	}

	return nil
}

func (vm *VirtualMachine) executeOp0xA() error {
	vm.I = vm.decodeNNN()
	return nil
}

func (vm *VirtualMachine) executeOp0xB() error {
	addr := (vm.decodeNNN() + vm.V[0x0]) & 0xfff
	if addr < ProgramStartAddress {
		return InvalidJumpError(vm.PC, addr)
	}

	vm.PC = addr
	return nil
}

func (vm *VirtualMachine) executeOp0xC() error {
	vm.V[vm.decodeX()] = uint(rand.Int()&0xff) & vm.decodeKK()
	return nil
}

func (vm *VirtualMachine) executeOp0xD() error {
	vm.V[0xf] = 0x0

	vx := vm.V[vm.decodeX()]
	vy := vm.V[vm.decodeY()]

	for n := uint(0); n < vm.decodeN(); n++ {
		y := (vy + n) % DisplayHeight
		sprite := vm.Memory[(vm.I+n)%MemorySize]

		for i := 7; sprite > 0x00; i-- {
			x := (vx + uint(i)) % DisplayWidth

			bit := sprite&0x1 == 0x1
			if bit && vm.Display[y][x] {
				vm.V[0xf] = 0x1
			}

			sprite >>= 1
			vm.Display[y][x] = vm.Display[y][x] != bit
		}
	}

	return nil
}

func (vm *VirtualMachine) executeOp0xE() error {
	vx := vm.V[vm.decodeX()]

	switch vm.decodeKK() {
	case 0x9e:
		if vm.Keys[vx] {
			vm.PC += 0x2
		}
	case 0xa1:
		if !vm.Keys[vx] {
			vm.PC += 0x2
		}
	default:
		return InvalidOpcodeError(vm.Opcode)
	}

	return nil
}

func (vm *VirtualMachine) executeOp0xF() error {
	x := vm.decodeX()

	switch vm.decodeKK() {
	case 0x07:
		vm.V[x] = vm.DT
	case 0x0a:
		for i, k := range vm.Keys {
			if k {
				vm.V[x] = uint(i)
				return nil
			}
		}
		vm.PC -= 0x2
	case 0x15:
		vm.DT = vm.V[x]
	case 0x18:
		vm.ST = vm.V[x]
	case 0x1E:
		vm.I = (vm.I + vm.V[x]) & 0xfff
	case 0x29:
		vm.I = vm.V[x] * FontSize
	case 0x33:
		vm.Memory[vm.I] = vm.V[x] / 100
		vm.Memory[vm.I+1] = (vm.V[x] % 100) / 10
		vm.Memory[vm.I+2] = vm.V[x] % 10
	case 0x55:
		for i := uint(0); i <= x; i++ {
			vm.Memory[vm.I+i] = vm.V[i]
		}
	case 0x65:
		for i := uint(0); i <= x; i++ {
			vm.V[i] = vm.Memory[vm.I+i]
		}
	}

	return nil
}

//=====================================================================
// Misc.
//=====================================================================

// PrintState prints the state of the virtual machine.
//
// This function exists primarily for debugging purposes.
func (vm *VirtualMachine) PrintState() {
	fmt.Println("----------------")

	// Display special registers
	fmt.Printf("I:         0x%.3X\n", vm.I)
	fmt.Printf("SP:        0x%.1X\n", vm.SP)
	fmt.Printf("PC:        0x%.3X\n", vm.PC)
	fmt.Printf("Delay:     0x%.2X\n", vm.DT)
	fmt.Printf("Sound:     0x%.2X\n", vm.ST)

	// Display call stack
	if vm.SP > 0x0 {
		fmt.Printf("\nStack:\n")

		for i := int(vm.SP) - 1; i >= 0; i-- {
			fmt.Printf("  0x%.1X: |0x%.3X|\n", i, vm.Stack[i])
		}
	}

	// Display general-purpose registers
	fmt.Printf("\nRegisters:")
	for i := 0; i < 4; i++ {
		fmt.Printf("\n  ")

		for j := 0; j < 4; j++ {
			fmt.Printf("|0x%.1X: 0x%.2X|", j*4+i, vm.V[j*4+i])
		}
	}

	fmt.Println("\n----------------")
}
