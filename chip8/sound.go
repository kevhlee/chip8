package chip8

type Sound struct {
	value uint8
}

func NewSound() *Sound {
	return &Sound{}
}

func (s *Sound) Write(value uint8) {
	s.value = value
}

func (s *Sound) Step() {
	if s.value > 0 {
		s.value--
	}
}
