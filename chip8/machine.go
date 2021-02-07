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

// Machine is the CHIP-8 virtual machine.
type Machine struct {
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

// NewMachine creates a CHIP-8 new virtual machine.
func NewMachine(keyboard Keyboard, screen Screen) *Machine {
	m := &Machine{
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

// LoadProgram loads a program into memory.
func (m *Machine) LoadProgram(bytes ...byte) error {
	if len(bytes) >= (MemorySize - ProgramStartAddress) {
		return fmt.Errorf("Program has too many bytes")
	}

	for i, b := range bytes {
		m.memory[ProgramStartAddress+i] = b
	}
	return nil
}

// Reset resets the virtual machine without erasing the current program.
func (m *Machine) Reset() {
	m.i = 0
	m.pc = ProgramStartAddress
	m.sp = 0
	m.delay = 0
	m.sound = 0

	for i := 0; i < len(m.registers); i++ {
		m.registers[i] = 0
	}

	for i := 0; i < len(m.stack); i++ {
		m.registers[i] = 0
	}

	m.Screen.ClearScreen()
}

// RunCycle executes a single CPU cycle.
func (m *Machine) RunCycle() error {
	opcode := m.fetchOp()
	return m.executeOpMap[opcode>>12](opcode)
}

// UpdateTimers updates the virtual machine's timers.
func (m *Machine) UpdateTimers() {
	if m.delay > 0 {
		m.delay--
	}
	if m.sound > 0 {
		m.sound--
	}
}

func (m *Machine) fetchOp() uint16 {
	opcode := (uint16(m.memory[m.pc]) << 8) | uint16(m.memory[m.pc+1])
	m.pc += 2
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

func (m *Machine) executeOp0(opcode uint16) error {
	switch opcode {
	// 00E0 - CLS
	case 0x00E0:
		m.Screen.ClearScreen()

	// 00EE - RET
	case 0x00EE:
		if m.sp < 1 {
			return fmt.Errorf("Empty call stack")
		}

		m.sp--
		m.pc = m.stack[m.sp]
	}

	return nil
}

func (m *Machine) executeOp1(opcode uint16) error {
	// 1nnn - JP addr
	m.pc = decodeAddr(opcode)
	return nil
}

func (m *Machine) executeOp2(opcode uint16) error {
	// 2nnn - CALL addr
	if m.sp > 0xF {
		return fmt.Errorf("Call stack overflow")
	}

	m.stack[m.sp] = m.pc
	m.sp++
	m.pc = decodeAddr(opcode)
	return nil
}

func (m *Machine) executeOp3(opcode uint16) error {
	// 3xkk - SE Vx, byte
	if m.registers[decodeX(opcode)] == decodeByte(opcode) {
		m.pc += 2
	}
	return nil
}

func (m *Machine) executeOp4(opcode uint16) error {
	// 4xkk - SNE Vx, byte
	if m.registers[decodeX(opcode)] != decodeByte(opcode) {
		m.pc += 2
	}
	return nil
}

func (m *Machine) executeOp5(opcode uint16) error {
	// 5xy0 - SE Vx, Vy
	if decodeNibb(opcode) == 0 && m.registers[decodeX(opcode)] == m.registers[decodeY(opcode)] {
		m.pc += 2
	}
	return nil
}

func (m *Machine) executeOp6(opcode uint16) error {
	// 6xkk - LD Vx, byte
	m.registers[decodeX(opcode)] = decodeByte(opcode)
	return nil
}

func (m *Machine) executeOp7(opcode uint16) error {
	// 7xkk - ADD Vx, byte
	m.registers[decodeX(opcode)] += decodeByte(opcode)
	return nil
}

func (m *Machine) executeOp8(opcode uint16) error {
	x, y := decodeX(opcode), decodeY(opcode)

	switch decodeNibb(opcode) {
	// 8xy0 - LD Vx, Vy
	case 0x0:
		m.registers[x] = m.registers[y]

	// 8xy1 - OR Vx, Vy
	case 0x1:
		m.registers[x] |= m.registers[y]

	// 8xy2 - AND Vx, Vy
	case 0x2:
		m.registers[x] &= m.registers[y]

	// 8xy3 - XOR Vx, Vy
	case 0x3:
		m.registers[x] ^= m.registers[y]

	// 8xy4 - ADD Vx, Vy
	case 0x4:
		if m.registers[x] > 0xFF-m.registers[y] {
			m.registers[0xF] = 1
		} else {
			m.registers[0xF] = 0
		}
		m.registers[x] += m.registers[y]

	// 8xy5 - SUB Vx, Vy
	case 0x5:
		if m.registers[x] > m.registers[y] {
			m.registers[0xF] = 1
		} else {
			m.registers[0xF] = 0
		}
		m.registers[x] -= m.registers[y]

	// 8xy6 - SHR Vx, {Vy}
	case 0x6:
		m.registers[0xF] = m.registers[x] & 1
		m.registers[x] >>= 1

	// 8xy7 - SUBN Vx, Vy
	case 0x7:
		if m.registers[y] > m.registers[x] {
			m.registers[0xF] = 1
		} else {
			m.registers[0xF] = 0
		}
		m.registers[x] = m.registers[y] - m.registers[x]

	// 8xyE - SHL Vx, {Vy}
	case 0xE:
		m.registers[0xF] = m.registers[x] >> 7
		m.registers[x] <<= 1
	}

	return nil
}

func (m *Machine) executeOp9(opcode uint16) error {
	// 9xy0 - SNE Vx, Vy
	if decodeNibb(opcode) == 0 && m.registers[decodeX(opcode)] != m.registers[decodeY(opcode)] {
		m.pc += 2
	}
	return nil
}

func (m *Machine) executeOpA(opcode uint16) error {
	// Annn - LD I, addr
	m.i = decodeAddr(opcode)
	return nil
}

func (m *Machine) executeOpB(opcode uint16) error {
	// Bnnn - JP addr, V0
	m.pc = decodeAddr(opcode) + uint16(m.registers[0])
	return nil
}

func (m *Machine) executeOpC(opcode uint16) error {
	// Cxkk - RND Vx, byte
	m.registers[decodeX(opcode)] = uint8(rand.Intn(0x100)) & decodeByte(opcode)
	return nil
}

func (m *Machine) executeOpD(opcode uint16) error {
	// Dxyn - DRW Vx, Vy, nibb
	vx := m.registers[decodeX(opcode)]
	vy := m.registers[decodeY(opcode)]

	if m.Screen.SetSprite(vx, vy, m.memory[m.i:m.i+uint16(decodeNibb(opcode))]...) {
		m.registers[0xF] = 1
	} else {
		m.registers[0xF] = 0
	}

	return nil
}

func (m *Machine) executeOpE(opcode uint16) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Ex9E - SKP Vx
	case 0x9E:
		if m.Keyboard.IsKeyPressed(m.registers[x]) {
			m.pc += 2
		}

	// ExA1 - SKNP Vx
	case 0xA1:
		if !m.Keyboard.IsKeyPressed(m.registers[x]) {
			m.pc += 2
		}
	}

	return nil
}

func (m *Machine) executeOpF(opcode uint16) error {
	x := decodeX(opcode)

	switch decodeByte(opcode) {
	// Fx07 - LD Vx, DT
	case 0x07:
		m.registers[x] = m.delay

	// Fx0A - LD Vx, K
	case 0x0A:
		if key, ok := m.Keyboard.PollKeyPress(); ok {
			m.registers[x] = key
		} else {
			m.pc -= 2
		}

	// Fx15 - LD DT, Vx
	case 0x15:
		m.delay = m.registers[x]

	// Fx18 - LD ST, Vx
	case 0x18:
		m.sound = m.registers[x]

	// Fx1E - ADD I, Vx
	case 0x1E:
		m.i += uint16(m.registers[x])

	// Fx29 - LD F, Vx
	case 0x29:
		m.i = uint16(m.registers[x]) * FontSize

	// Fx33 - LD B, Vx
	case 0x33:
		m.memory[m.i] = m.registers[x] / 100
		m.memory[m.i+1] = (m.registers[x] % 100) / 10
		m.memory[m.i+2] = m.registers[x] % 10

	// Fx55 - LD [I], Vx
	case 0x55:
		for i := uint16(0); i <= uint16(x); i++ {
			m.memory[m.i+i] = m.registers[i]
		}

	// Fx65 - LD Vx, [I]
	case 0x65:
		for i := uint16(0); i <= uint16(x); i++ {
			m.registers[i] = m.memory[m.i+i]
		}
	}

	return nil
}
