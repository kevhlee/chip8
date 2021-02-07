package chip8

const (
	// NumKeys is the number of keys in the CHIP-8 hexadecimal keyboard.
	NumKeys = 0x10
)

// Keyboard is the CHIP-8 hexadecimal keyboard.
type Keyboard interface {
	IsKeyPressed(key uint8) bool
	PollKeyPress() (uint8, bool)
}
