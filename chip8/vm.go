package chip8

import (
	"fmt"
	"math/rand"
)

const (
	// ProgramStartAddress is the start memory address for CHIP-8 programs.
	ProgramStartAddress = 0x200
)

// VirtualMachine is the CHIP-8 virtual machine.
type VirtualMachine struct {
	i      uint16
	sp     uint8
	pc     uint16
	memory [0x1000]uint8
	v      [0x10]uint8
	stack  [0x10]uint16
}

// NewVirtualMachine creates a CHIP-8 new virtual machine.
func NewVirtualMachine() *VirtualMachine {
	m := &VirtualMachine{
		pc:    ProgramStartAddress,
		v:     [0x10]uint8{},
		stack: [0x10]uint16{},
	}

	m.memory = [0x1000]uint8{
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

	return m
}

// LoadProgram loads bytes into memory.
func (vm *VirtualMachine) LoadProgram(bytes ...byte) error {
	if len(bytes) >= (len(vm.memory) - ProgramStartAddress) {
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

	for i := 0; i < len(vm.v); i++ {
		vm.v[i] = 0
	}

	for i := 0; i < len(vm.stack); i++ {
		vm.v[i] = 0
	}
}

// Step executes a single CPU cycle.
func (vm *VirtualMachine) Step(keyboard *Keyboard, screen *Screen, sound *Sound, timer *Timer) error {
	opcode := vm.fetchOpcode()

	switch opcode >> 12 {
	case 0x0:
		return vm.executeOp0(opcode, screen)
	case 0x1:
		return vm.executeOp1(opcode)
	case 0x2:
		return vm.executeOp2(opcode)
	case 0x3:
		return vm.executeOp3(opcode)
	case 0x4:
		return vm.executeOp4(opcode)
	case 0x5:
		return vm.executeOp5(opcode)
	case 0x6:
		return vm.executeOp6(opcode)
	case 0x7:
		return vm.executeOp7(opcode)
	case 0x8:
		return vm.executeOp8(opcode)
	case 0x9:
		return vm.executeOp9(opcode)
	case 0xA:
		return vm.executeOpA(opcode)
	case 0xB:
		return vm.executeOpB(opcode)
	case 0xC:
		return vm.executeOpC(opcode)
	case 0xD:
		return vm.executeOpD(opcode, screen)
	case 0xE:
		return vm.executeOpE(opcode, keyboard)
	case 0xF:
		return vm.executeOpF(opcode, keyboard, sound, timer)
	default:
		// Unreachable
		return nil
	}
}

func (vm *VirtualMachine) fetchOpcode() uint16 {
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

func (vm *VirtualMachine) executeOp0(opcode uint16, screen *Screen) error {
	switch opcode {
	// 00E0 - CLS
	case 0x00E0:
		screen.Clear()

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
	if vm.v[decodeX(opcode)] == decodeByte(opcode) {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp4(opcode uint16) error {
	// 4xkk - SNE Vx, byte
	if vm.v[decodeX(opcode)] != decodeByte(opcode) {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp5(opcode uint16) error {
	// 5xy0 - SE Vx, Vy
	if decodeNibb(opcode) == 0 && vm.v[decodeX(opcode)] == vm.v[decodeY(opcode)] {
		vm.pc += 2
	}
	return nil
}

func (vm *VirtualMachine) executeOp6(opcode uint16) error {
	// 6xkk - LD Vx, byte
	vm.v[decodeX(opcode)] = decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp7(opcode uint16) error {
	// 7xkk - ADD Vx, byte
	vm.v[decodeX(opcode)] += decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOp8(opcode uint16) error {
	x, y := decodeX(opcode), decodeY(opcode)

	switch decodeNibb(opcode) {
	// 8xy0 - LD Vx, Vy
	case 0x0:
		vm.v[x] = vm.v[y]

	// 8xy1 - OR Vx, Vy
	case 0x1:
		vm.v[x] |= vm.v[y]

	// 8xy2 - AND Vx, Vy
	case 0x2:
		vm.v[x] &= vm.v[y]

	// 8xy3 - XOR Vx, Vy
	case 0x3:
		vm.v[x] ^= vm.v[y]

	// 8xy4 - ADD Vx, Vy
	case 0x4:
		if vm.v[x] > 0xFF-vm.v[y] {
			vm.v[0xF] = 1
		} else {
			vm.v[0xF] = 0
		}
		vm.v[x] += vm.v[y]

	// 8xy5 - SUB Vx, Vy
	case 0x5:
		if vm.v[x] > vm.v[y] {
			vm.v[0xF] = 1
		} else {
			vm.v[0xF] = 0
		}
		vm.v[x] -= vm.v[y]

	// 8xy6 - SHR Vx, {Vy}
	case 0x6:
		vm.v[0xF] = vm.v[x] & 1
		vm.v[x] >>= 1

	// 8xy7 - SUBN Vx, Vy
	case 0x7:
		if vm.v[y] > vm.v[x] {
			vm.v[0xF] = 1
		} else {
			vm.v[0xF] = 0
		}
		vm.v[x] = vm.v[y] - vm.v[x]

	// 8xyE - SHL Vx, {Vy}
	case 0xE:
		vm.v[0xF] = vm.v[x] >> 7
		vm.v[x] <<= 1
	}

	return nil
}

func (vm *VirtualMachine) executeOp9(opcode uint16) error {
	// 9xy0 - SNE Vx, Vy
	if decodeNibb(opcode) == 0 && vm.v[decodeX(opcode)] != vm.v[decodeY(opcode)] {
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
	vm.pc = decodeAddr(opcode) + uint16(vm.v[0])
	return nil
}

func (vm *VirtualMachine) executeOpC(opcode uint16) error {
	// Cxkk - RND Vx, byte
	vm.v[decodeX(opcode)] = uint8(rand.Intn(0x100)) & decodeByte(opcode)
	return nil
}

func (vm *VirtualMachine) executeOpD(opcode uint16, screen *Screen) error {
	// Dxyn - DRW Vx, Vy, nibb
	vx := vm.v[decodeX(opcode)]
	vy := vm.v[decodeY(opcode)]

	if screen.SetSprite(vx, vy, vm.memory[vm.i:vm.i+uint16(decodeNibb(opcode))]...) {
		vm.v[0xF] = 1
	} else {
		vm.v[0xF] = 0
	}

	return nil
}

func (vm *VirtualMachine) executeOpE(opcode uint16, keyboard *Keyboard) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Ex9E - SKP Vx
	case 0x9E:
		if keyboard.IsPressed(vm.v[x]) {
			vm.pc += 2
		}

	// ExA1 - SKNP Vx
	case 0xA1:
		if !keyboard.IsPressed(vm.v[x]) {
			vm.pc += 2
		}
	}

	return nil
}

func (vm *VirtualMachine) executeOpF(opcode uint16, keyboard *Keyboard, sound *Sound, timer *Timer) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Fx07 - LD Vx, DT
	case 0x07:
		vm.v[x] = timer.Read()

	// Fx0A - LD Vx, K
	case 0x0A:
		if key, ok := keyboard.Poll(); ok {
			vm.v[x] = key
		} else {
			vm.pc -= 2
		}

	// Fx15 - LD DT, Vx
	case 0x15:
		timer.Write(vm.v[x])

	// Fx18 - LD ST, Vx
	case 0x18:
		sound.Write(vm.v[x])

	// Fx1E - ADD I, Vx
	case 0x1E:
		vm.i += uint16(vm.v[x])

	// Fx29 - LD F, Vx
	case 0x29:
		vm.i = uint16(vm.v[x]) * 0x5

	// Fx33 - LD B, Vx
	case 0x33:
		vm.memory[vm.i] = vm.v[x] / 100
		vm.memory[vm.i+1] = (vm.v[x] % 100) / 10
		vm.memory[vm.i+2] = vm.v[x] % 10

	// Fx55 - LD [I], Vx
	case 0x55:
		for i := uint16(0); i <= uint16(x); i++ {
			vm.memory[vm.i+i] = vm.v[i]
		}

	// Fx65 - LD Vx, [I]
	case 0x65:
		for i := uint16(0); i <= uint16(x); i++ {
			vm.v[i] = vm.memory[vm.i+i]
		}
	}

	return nil
}
