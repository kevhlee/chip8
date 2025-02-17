package chip8

const (
	// NumKeys is the number of keys in the CHIP-8 hexadecimal keyboard.
	NumKeys = 0x10
)

// Keyboard is the CHIP-8 hexadecimal keyboard.
type Keyboard struct {
	keys    [NumKeys]bool
	polling bool
	lastKey int
}

func NewKeyboard() *Keyboard {
	return &Keyboard{
		keys:    [NumKeys]bool{},
		polling: false,
		lastKey: -1,
	}
}

func (k *Keyboard) IsPressed(key uint8) bool {
	return k.keys[key]
}

func (k *Keyboard) Set(key uint8, pressed bool) {
	if k.polling && k.lastKey == -1 {
		k.lastKey = int(key)
	}
	k.keys[key] = pressed
}

func (k *Keyboard) Poll() (uint8, bool) {
	if k.polling && k.lastKey >= 0 {
		key := uint8(k.lastKey)
		k.polling = false
		k.lastKey = -1
		return key, true
	} else {
		k.polling = true
		return 0, false
	}
}
