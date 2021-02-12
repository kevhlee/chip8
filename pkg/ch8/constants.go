package ch8

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

	// DefaultScale is the default scale factor of the CHIP-8 screen.
	DefaultScale = 10

	// DefaultFrequency is the default frequency of the CHIP-8 beeper.
	DefaultFrequency = 440

	// DefaultSampleRate is the default sample rate of the CHIP-8
	// beeper.
	DefaultSampleRate = 44100
)
