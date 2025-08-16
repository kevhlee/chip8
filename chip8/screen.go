package chip8

import "github.com/veandco/go-sdl2/sdl"

const (
	// ScreenWidth is the width of the CHIP-8 screen.
	ScreenWidth = 0x40

	// ScreenHeight is the height of the CHIP-8 screen.
	ScreenHeight = 0x20

	// ScreenSize is the number of pixels in the CHIP-8 screen.
	ScreenSize = ScreenWidth * ScreenHeight
)

// Screen is the CHIP-8 screen.
type Screen struct {
	buffer   [ScreenSize]bool
	rect     *sdl.Rect
	renderer *sdl.Renderer
}

func NewScreen(renderer *sdl.Renderer, scale int32) *Screen {
	return &Screen{
		buffer:   [ScreenSize]bool{},
		rect:     &sdl.Rect{W: scale, H: scale},
		renderer: renderer,
	}
}

func (s *Screen) Render() {
	s.renderer.SetDrawColor(0x00, 0x00, 0x00, 0xFF)
	s.renderer.Clear()
	s.renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)

	for i, on := range s.buffer {
		if on {
			s.rect.X = int32(i%ScreenWidth) * s.rect.W
			s.rect.Y = int32(i/ScreenWidth) * s.rect.H
			s.renderer.FillRect(s.rect)
		}
	}

	s.renderer.Present()
}

func (s *Screen) Clear() {
	for i := 0; i < len(s.buffer); i++ {
		s.buffer[i] = false
	}
}

func (s *Screen) SetSprite(x, y uint8, sprite ...uint8) bool {
	flag := false

	for i, b := range sprite {
		for j := uint8(0); j < 8; j++ {
			idx := s.getBufferIndex(x+(7-j), y+uint8(i))
			bit := b&1 == 1

			if s.buffer[idx] && bit {
				flag = true
			}

			s.buffer[idx] = s.buffer[idx] != bit
			b >>= 1
		}
	}

	return flag
}

func (s *Screen) getBufferIndex(x, y uint8) int {
	return (int(y%ScreenHeight) * ScreenWidth) + int(x%ScreenWidth)
}
