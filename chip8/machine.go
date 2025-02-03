package chip8

import (
	"fmt"
	"math/rand"
)

const (
	// FontSize is the number of bytes in of a CHIP-8 hexadecimal sprite.
	FontSize = 0x5

	// MemorySize is the total size of the virtual machine's memory.
	MemorySize = 0x1000

	// NumRegisters is the number of general-purpose registers (V) in the
	// virtual machine.
	NumRegisters = 0x10

	// ProgramStartAddress is the start memory address for CHIP-8 programs.
	ProgramStartAddress = 0x200

	// StackSize is the total size of the virtual machine's call stack.
	StackSize = 0x10
)

// VirtualMachine is the CHIP-8 virtual machine.
type VirtualMachine struct {
	Keyboard
	Screen

	i            uint16
	sp           uint8
	pc           uint16
	delay        uint8
	sound        uint8
	memory       [MemorySize]uint8
	registers    [NumRegisters]uint8
	stack        [StackSize]uint16
	executeOpMap [0x10]func(uint16) error
}

// NewVirtualMachine creates a CHIP-8 new virtual machine.
func NewVirtualMachine(keyboard Keyboard, screen Screen) *VirtualMachine {
	m := &VirtualMachine{
		pc:        ProgramStartAddress,
		registers: [NumRegisters]uint8{},
		stack:     [StackSize]uint16{},
		Keyboard:  keyboard,
		Screen:    screen,
	}

	m.memory = [MemorySize]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	m.executeOpMap = [0x10]func(uint16) error{
		m.executeOp0, m.executeOp1, m.executeOp2, m.executeOp3,
		m.executeOp4, m.executeOp5, m.executeOp6, m.executeOp7,
		m.executeOp8, m.executeOp9, m.executeOpA, m.executeOpB,
		m.executeOpC, m.executeOpD, m.executeOpE, m.executeOpF,
	}

	return m
}

// LoadProgram loads bytes into memory.
func (vm *VirtualMachine) LoadProgram(bytes ...byte) error {
	if len(bytes) >= (MemorySize - ProgramStartAddress) {
		return fmt.Errorf("Program has too many bytes")
	}

	for i, b := range bytes {
		vm.memory[ProgramStartAddress+i] = b
	}
	return nil
}

// Reset resets the virtual machine without erasing the current program.
func (vm *VirtualMachine) Reset() {
	vm.i = 0
	vm.pc = ProgramStartAddress
	vm.sp = 0
	vm.delay = 0
	vm.sound = 0

	for i := 0; i < len(vm.registers); i++ {
		vm.registers[i] = 0
	}

	for i := 0; i < len(vm.stack); i++ {
		vm.registers[i] = 0
	}

	vm.Screen.ClearScreen()
}

// RunCycle executes a single CPU cycle.
func (vm *VirtualMachine) RunCycle() error {
	opcode := vm.fetchOp()
	return vm.executeOpMap[opcode>>12](opcode)
}

// UpdateTimers updates the virtual machine's timers.
func (vm *VirtualMachine) UpdateTimers() {
	if vm.delay > 0 {
		vm.delay--
	}
	if vm.sound > 0 {
		vm.sound--
	}
}

func (vm *VirtualMachine) fetchOp() uint16 {
	opcode := (uint16(vm.memory[vm.pc]) << 8) | uint16(vm.memory[vm.pc+1])
	vm.pc += 2
	return opcode
}

func decodeX(opcode uint16) uint8 {
	return uint8((opcode >> 8) & 0xF)
}

func decodeY(opcode uint16) uint8 {
	return uint8((opcode >> 4) & 0xF)
}

func decodeAddr(opcode uint16) uint16 {
	return opcode & 0xFFF
}

func decodeByte(opcode uint16) uint8 {
	return uint8(opcode & 0xFF)
}

func decodeNibb(opcode uint16) uint8 {
	return uint8(opcode & 0xF)
}

func (vm *VirtualMachine) executeOp0(opcode uint16) error {
	switch opcode {
	// 00E0 - CLS
	case 0x00E0:
		vm.Screen.ClearScreen()

	// 00EE - RET
	case 0x00EE:
		if vm.sp < 1 {
			return fmt.Errorf("Empty call stack")
		}

		vm.sp--
		vm.pc = vm.stack[vm.sp]
	}

	return nil
}

func (vm *VirtualMachine) executeOp1(opcode uint16) error {
	// 1nnn - JP addr
	vm.pc = decodeAddr(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp2(opcode uint16) error {
	// 2nnn - CALL addr
	if vm.sp > 0xF {
		return fmt.Errorf("Call stack overflow")
	}

	vm.stack[vm.sp] = vm.pc
	vm.sp++
	vm.pc = decodeAddr(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp3(opcode uint16) error {
	// 3xkk - SE Vx, byte
	if vm.registers[decodeX(opcode)] == decodeByte(opcode) {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp4(opcode uint16) error {
	// 4xkk - SNE Vx, byte
	if vm.registers[decodeX(opcode)] != decodeByte(opcode) {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp5(opcode uint16) error {
	// 5xy0 - SE Vx, Vy
	if decodeNibb(opcode) == 0 && vm.registers[decodeX(opcode)] == vm.registers[decodeY(opcode)] {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp6(opcode uint16) error {
	// 6xkk - LD Vx, byte
	vm.registers[decodeX(opcode)] = decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp7(opcode uint16) error {
	// 7xkk - ADD Vx, byte
	vm.registers[decodeX(opcode)] += decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp8(opcode uint16) error {
	x, y := decodeX(opcode), decodeY(opcode)

	switch decodeNibb(opcode) {
	// 8xy0 - LD Vx, Vy
	case 0x0:
		vm.registers[x] = vm.registers[y]

	// 8xy1 - OR Vx, Vy
	case 0x1:
		vm.registers[x] |= vm.registers[y]

	// 8xy2 - AND Vx, Vy
	case 0x2:
		vm.registers[x] &= vm.registers[y]

	// 8xy3 - XOR Vx, Vy
	case 0x3:
		vm.registers[x] ^= vm.registers[y]

	// 8xy4 - ADD Vx, Vy
	case 0x4:
		if vm.registers[x] > 0xFF-vm.registers[y] {
			vm.registers[0xF] = 1
		} else {
			vm.registers[0xF] = 0
		}
		vm.registers[x] += vm.registers[y]

	// 8xy5 - SUB Vx, Vy
	case 0x5:
		if vm.registers[x] > vm.registers[y] {
			vm.registers[0xF] = 1
		} else {
			vm.registers[0xF] = 0
		}
		vm.registers[x] -= vm.registers[y]

	// 8xy6 - SHR Vx, {Vy}
	case 0x6:
		vm.registers[0xF] = vm.registers[x] & 1
		vm.registers[x] >>= 1

	// 8xy7 - SUBN Vx, Vy
	case 0x7:
		if vm.registers[y] > vm.registers[x] {
			vm.registers[0xF] = 1
		} else {
			vm.registers[0xF] = 0
		}
		vm.registers[x] = vm.registers[y] - vm.registers[x]

	// 8xyE - SHL Vx, {Vy}
	case 0xE:
		vm.registers[0xF] = vm.registers[x] >> 7
		vm.registers[x] <<= 1
	}

	return nil
}

func (vm *VirtualMachine) executeOp9(opcode uint16) error {
	// 9xy0 - SNE Vx, Vy
	if decodeNibb(opcode) == 0 && vm.registers[decodeX(opcode)] != vm.registers[decodeY(opcode)] {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOpA(opcode uint16) error {
	// Annn - LD I, addr
	vm.i = decodeAddr(opcode)
	return nil
}

func (vm *VirtualMachine) executeOpB(opcode uint16) error {
	// Bnnn - JP addr, V0
	vm.pc = decodeAddr(opcode) + uint16(vm.registers[0])
	return nil
}

func (vm *VirtualMachine) executeOpC(opcode uint16) error {
	// Cxkk - RND Vx, byte
	vm.registers[decodeX(opcode)] = uint8(rand.Intn(0x100)) & decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOpD(opcode uint16) error {
	// Dxyn - DRW Vx, Vy, nibb
	vx := vm.registers[decodeX(opcode)]
	vy := vm.registers[decodeY(opcode)]

	if vm.Screen.SetSprite(vx, vy, vm.memory[vm.i:vm.i+uint16(decodeNibb(opcode))]...) {
		vm.registers[0xF] = 1
	} else {
		vm.registers[0xF] = 0
	}

	return nil
}

func (vm *VirtualMachine) executeOpE(opcode uint16) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Ex9E - SKP Vx
	case 0x9E:
		if vm.Keyboard.IsKeyPressed(vm.registers[x]) {
			vm.pc += 2
		}

	// ExA1 - SKNP Vx
	case 0xA1:
		if !vm.Keyboard.IsKeyPressed(vm.registers[x]) {
			vm.pc += 2
		}
	}

	return nil
}

func (vm *VirtualMachine) executeOpF(opcode uint16) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Fx07 - LD Vx, DT
	case 0x07:
		vm.registers[x] = vm.delay

	// Fx0A - LD Vx, K
	case 0x0A:
		if key, ok := vm.Keyboard.PollKeyPress(); ok {
			vm.registers[x] = key
		} else {
			vm.pc -= 2
		}

	// Fx15 - LD DT, Vx
	case 0x15:
		vm.delay = vm.registers[x]

	// Fx18 - LD ST, Vx
	case 0x18:
		vm.sound = vm.registers[x]

	// Fx1E - ADD I, Vx
	case 0x1E:
		vm.i += uint16(vm.registers[x])

	// Fx29 - LD F, Vx
	case 0x29:
		vm.i = uint16(vm.registers[x]) * FontSize

	// Fx33 - LD B, Vx
	case 0x33:
		vm.memory[vm.i] = vm.registers[x] / 100
		vm.memory[vm.i+1] = (vm.registers[x] % 100) / 10
		vm.memory[vm.i+2] = vm.registers[x] % 10

	// Fx55 - LD [I], Vx
	case 0x55:
		for i := uint16(0); i <= uint16(x); i++ {
			vm.memory[vm.i+i] = vm.registers[i]
		}

	// Fx65 - LD Vx, [I]
	case 0x65:
		for i := uint16(0); i <= uint16(x); i++ {
			vm.registers[i] = vm.memory[vm.i+i]
		}
	}

	return nil
}
