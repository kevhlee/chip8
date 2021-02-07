package ch8

import (
	"fmt"
	"io/ioutil"
	"math/rand"
)

var (
	fonts = []uint{
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
)

//===========================================================================
// Virtual Machine
//===========================================================================

// VirtualMachine is the CHIP-8 virtual machine.
type VirtualMachine struct {
	I           uint
	SP          uint
	PC          uint
	Delay       uint
	Sound       uint
	Stack       [StackSize]uint
	Memory      [MemorySize]uint
	Registers   [NumberOfRegisters]uint
	Keys        [NumberOfKeys]bool
	Display     [NumberOfPixels]uint
	OpcodeExecs map[uint]func(uint)
}

// NewVirtualMachine creates new CHIP-8 virtual machine instance.
func NewVirtualMachine() *VirtualMachine {
	vm := &VirtualMachine{
		PC:        ProgramStartAddress,
		Stack:     [StackSize]uint{},
		Registers: [NumberOfRegisters]uint{},
		Keys:      [NumberOfKeys]bool{},
		Display:   [NumberOfPixels]uint{},
		Memory:    [MemorySize]uint{},
	}

	for i, b := range fonts {
		vm.Memory[i] = b
	}

	vm.OpcodeExecs = map[uint]func(uint){
		0x0: vm.executeOp0x0,
		0x1: vm.executeOp0x1,
		0x2: vm.executeOp0x2,
		0x3: vm.executeOp0x3,
		0x4: vm.executeOp0x4,
		0x5: vm.executeOp0x5,
		0x6: vm.executeOp0x6,
		0x7: vm.executeOp0x7,
		0x8: vm.executeOp0x8,
		0x9: vm.executeOp0x9,
		0xa: vm.executeOp0xA,
		0xb: vm.executeOp0xB,
		0xc: vm.executeOp0xC,
		0xd: vm.executeOp0xD,
		0xe: vm.executeOp0xE,
		0xf: vm.executeOp0xF,
	}

	return vm
}

// RunCycle runs a single CPU cycle of the virtual machine.
func (vm *VirtualMachine) RunCycle() {
	// Fetch next opcode
	opcode := vm.fetch()

	// Decode and execute opcode
	execute := vm.decode(opcode)

	// Execute opcode
	execute(opcode)

	// Keep program counter within range
	if vm.PC > 0xfff {
		vm.PC = (vm.PC & 0xfff) + ProgramStartAddress
	}
}

// UpdateTimers updates the delay and sound timers.
func (vm *VirtualMachine) UpdateTimers() {
	if vm.Delay > 0x00 {
		vm.Delay--
	}
	if vm.Sound > 0x00 {
		vm.Sound--
	}
}

// LoadROM reads a CHIP-8 ROM program file (*.ch8) and loads it into
// memory.
func (vm *VirtualMachine) LoadROM(path string) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	} else if len(data) > ProgramMemorySize {
		return fmt.Errorf("Not enough memory for opcodes")
	}

	i := ProgramStartAddress

	for _, b := range data {
		vm.Memory[i] = uint(b)
		i++
	}

	return nil
}

// LoadOpcodes loads opcodes into the virtual machine's program memory.
func (vm *VirtualMachine) LoadOpcodes(opcodes []uint) {
	i := ProgramStartAddress
	for _, opcode := range opcodes {
		vm.Memory[i] = opcode >> 8
		vm.Memory[i+1] = opcode & 0xff
		i += 0x2
	}
}

// Reset resets the entire state of the virtual machine.
func (vm *VirtualMachine) Reset() {
	vm.I = 0x000
	vm.SP = 0x00
	vm.PC = ProgramStartAddress
	vm.Delay = 0x00
	vm.Sound = 0x00

	vm.ResetDisplay()

	for i := 0; i < len(vm.Keys); i++ {
		vm.Keys[i] = false
	}

	for i := 0; i < len(vm.Stack); i++ {
		vm.Stack[i] = 0x000
	}

	for i := 0; i < len(vm.Registers); i++ {
		vm.Registers[i] = 0x00
	}

	for i := ProgramStartAddress; i < len(vm.Memory); i++ {
		vm.Memory[i] = 0x00
	}
}

// ResetDisplay resets the state of the display.
func (vm *VirtualMachine) ResetDisplay() {
	for i := 0; i < len(vm.Display); i++ {
		vm.Display[i] = 0x00
	}
}

func (vm *VirtualMachine) fetch() uint {
	opcode := (vm.Memory[vm.PC] << 8) | vm.Memory[vm.PC+1]
	vm.PC += 0x2
	return opcode
}

func (vm *VirtualMachine) decode(opcode uint) func(uint) {
	return vm.OpcodeExecs[getOp(opcode)]
}

//=====================================================================
// Decode
//=====================================================================

func getOp(opcode uint) uint {
	return opcode >> 0xc
}

func getX(opcode uint) uint {
	return (opcode >> 8) & 0xf
}

func getY(opcode uint) uint {
	return (opcode >> 4) & 0xf
}

func getN(opcode uint) uint {
	return opcode & 0xf
}

func getKK(opcode uint) uint {
	return opcode & 0xff
}

func getNNN(opcode uint) uint {
	return opcode & 0xfff
}

//=====================================================================
// Execute
//=====================================================================

func (vm *VirtualMachine) executeOp0x0(opcode uint) {
	switch getNNN(opcode) {
	case 0x0e0:
		vm.ResetDisplay()
	case 0x0ee:
		vm.SP--
		vm.PC = vm.Stack[vm.SP]
	}
}

func (vm *VirtualMachine) executeOp0x1(opcode uint) {
	vm.PC = getNNN(opcode)
}

func (vm *VirtualMachine) executeOp0x2(opcode uint) {
	vm.Stack[vm.SP] = vm.PC
	vm.PC = getNNN(opcode)
	vm.SP++
}

func (vm *VirtualMachine) executeOp0x3(opcode uint) {
	if vm.Registers[getX(opcode)] == getKK(opcode) {
		vm.PC += 2
	}
}

func (vm *VirtualMachine) executeOp0x4(opcode uint) {
	if vm.Registers[getX(opcode)] != getKK(opcode) {
		vm.PC += 2
	}
}

func (vm *VirtualMachine) executeOp0x5(opcode uint) {
	if getN(opcode) == 0x0 && vm.Registers[getX(opcode)] == getKK(opcode) {
		vm.PC += 2
	}
}

func (vm *VirtualMachine) executeOp0x6(opcode uint) {
	vm.Registers[getX(opcode)] = getKK(opcode)
}

func (vm *VirtualMachine) executeOp0x7(opcode uint) {
	x := getX(opcode)
	kk := getKK(opcode)
	vm.Registers[x] = (vm.Registers[x] + kk) & 0xff
}

func (vm *VirtualMachine) executeOp0x8(opcode uint) {
	x := getX(opcode)
	y := getY(opcode)

	switch getN(opcode) {
	case 0x0:
		vm.Registers[x] = vm.Registers[y]
	case 0x1:
		vm.Registers[x] |= vm.Registers[y]
	case 0x2:
		vm.Registers[x] &= vm.Registers[y]
	case 0x3:
		vm.Registers[x] ^= vm.Registers[y]
	case 0x4:
		result := vm.Registers[x] + vm.Registers[y]
		if result > 0xff {
			vm.Registers[0xf] = 0x1
			vm.Registers[x] = result & 0xff
		} else {
			vm.Registers[0xf] = 0x0
			vm.Registers[x] = result
		}
	case 0x5:
		if vm.Registers[x] > vm.Registers[y] {
			vm.Registers[0xf] = 0x1
		} else {
			vm.Registers[0xf] = 0x0
		}
		vm.Registers[x] = (vm.Registers[x] - vm.Registers[y]) & 0xff
	case 0x6:
		vm.Registers[0xf] = vm.Registers[x] & 0x01
		vm.Registers[x] >>= 1
	case 0x7:
		if vm.Registers[x] < vm.Registers[y] {
			vm.Registers[0xf] = 0x1
		} else {
			vm.Registers[0xf] = 0x0
		}
		vm.Registers[x] = (vm.Registers[y] - vm.Registers[x]) & 0xff
	case 0xe:
		vm.Registers[0xf] = vm.Registers[x] >> 7
		vm.Registers[x] = (vm.Registers[x] << 1) & 0xff
	default:
		panic(fmt.Sprintf("%.4X\n", opcode))
	}
}

func (vm *VirtualMachine) executeOp0x9(opcode uint) {
	if getN(opcode) == 0x0 && vm.Registers[getX(opcode)] != getKK(opcode) {
		vm.PC += 2
	}
}

func (vm *VirtualMachine) executeOp0xA(opcode uint) {
	vm.I = getNNN(opcode)
}

func (vm *VirtualMachine) executeOp0xB(opcode uint) {
	vm.PC = (getNNN(opcode) + vm.Registers[0x0]) & 0xfff
}

func (vm *VirtualMachine) executeOp0xC(opcode uint) {
	vm.Registers[getX(opcode)] = uint(rand.Int()&0xff) & getKK(opcode)
}

func (vm *VirtualMachine) executeOp0xD(opcode uint) {
	vm.Registers[0xf] = 0x0

	xStart := vm.Registers[getX(opcode)]
	yStart := vm.Registers[getY(opcode)]

	for n := uint(0); n < getN(opcode); n++ {
		y := (yStart + n) % DisplayHeight
		sprite := vm.Memory[(vm.I+n)%MemorySize]

		for i := 7; i >= 0; i-- {
			x := (xStart + uint(i)) % DisplayWidth
			idx := y*DisplayWidth + x

			bit := sprite & 0x1
			sprite >>= 1

			if bit&vm.Display[idx] == 0x1 {
				vm.Registers[0xf] = 0x1
			}

			vm.Display[idx] ^= bit
		}
	}
}

func (vm *VirtualMachine) executeOp0xE(opcode uint) {
	vx := vm.Registers[getX(opcode)]

	switch getKK(opcode) {
	case 0x9e:
		if vm.Keys[vx] {
			vm.PC += 0x2
		}
	case 0xa1:
		if !vm.Keys[vx] {
			vm.PC += 0x2
		}
	}
}

func (vm *VirtualMachine) executeOp0xF(opcode uint) {
	x := getX(opcode)

	switch getKK(opcode) {
	case 0x07:
		vm.Registers[x] = vm.Delay
	case 0x0a:
		for i, k := range vm.Keys {
			if k {
				vm.Registers[x] = uint(i)
				return
			}
		}
		vm.PC -= 0x2
	case 0x15:
		vm.Delay = vm.Registers[x]
	case 0x18:
		vm.Sound = vm.Registers[x]
	case 0x1E:
		vm.I = (vm.I + vm.Registers[x]) & 0xfff
	case 0x29:
		vm.I = vm.Registers[x] * FontSize
	case 0x33:
		vm.Memory[vm.I] = vm.Registers[x] / 100
		vm.Memory[vm.I+1] = (vm.Registers[x] % 100) / 10
		vm.Memory[vm.I+2] = vm.Registers[x] % 10
	case 0x55:
		for i := uint(0); i <= x; i++ {
			vm.Memory[vm.I+i] = vm.Registers[i]
		}
	case 0x65:
		for i := uint(0); i <= x; i++ {
			vm.Registers[i] = vm.Memory[vm.I+i]
		}
	}
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
	fmt.Printf("Delay:     0x%.2X\n", vm.Delay)
	fmt.Printf("Sound:     0x%.2X\n", vm.Sound)

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
			fmt.Printf("|0x%.1X: 0x%.2X|", j*4+i, vm.Registers[j*4+i])
		}
	}

	fmt.Println("\n----------------")
}
