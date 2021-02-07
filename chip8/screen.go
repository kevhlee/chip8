package chip8

const (
	// ScreenWidth is the width of the CHIP-8 screen.
	ScreenWidth = 0x40

	// ScreenHeight is the height of the CHIP-8 screen.
	ScreenHeight = 0x20

	// ScreenSize is the number of pixels in the CHIP-8 screen.
	ScreenSize = ScreenWidth * ScreenHeight
)

// Screen is the CHIP-8 screen.
type Screen interface {
	ClearScreen()
	SetSprite(x, y uint8, sprite ...uint8) bool
}
