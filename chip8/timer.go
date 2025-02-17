package chip8

type Timer struct {
	value uint8
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t Timer) Read() uint8 {
	return t.value
}

func (t *Timer) Write(value uint8) {
	t.value = value
}

func (t *Timer) Step() {
	if t.value > 0 {
		t.value--
	}
}
