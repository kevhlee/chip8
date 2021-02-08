package emu

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/kevhlee/chip8/pkg/ch8"
)

const (
	frequency  = 440
	sampleRate = 44100
)

// stream is an infinite stream of 440 Hz sine wave.
//
// This struct is taken directly from Ebiten's example code:
// https://ebiten.org/examples/sinewave.html
type stream struct {
	position  int64
	remaining []byte
}

// Read is io.Reader's Read.
//
// Read fills the data with sine wave samples.
func (s *stream) Read(buf []byte) (int, error) {
	if len(s.remaining) > 0 {
		n := copy(buf, s.remaining)
		s.remaining = s.remaining[n:]
		return n, nil
	}

	var origBuf []byte
	if len(buf)%4 > 0 {
		origBuf = buf
		buf = make([]byte, len(origBuf)+4-len(origBuf)%4)
	}

	const length = int64(sampleRate / frequency)
	p := s.position / 4
	for i := 0; i < len(buf)/4; i++ {
		const max = 32767
		b := int16(math.Sin(2*math.Pi*float64(p)/float64(length)) * max)
		buf[4*i] = byte(b)
		buf[4*i+1] = byte(b >> 8)
		buf[4*i+2] = byte(b)
		buf[4*i+3] = byte(b >> 8)
		p++
	}

	s.position += int64(len(buf))
	s.position %= length * 4

	if origBuf != nil {
		n := copy(origBuf, buf)
		s.remaining = buf[n:]
		return n, nil
	}
	return len(buf), nil
}

// Close is io.Closer's Close.
func (s *stream) Close() error {
	return nil
}

// Sound is the CHIP-8 emulator's audio.
type Sound struct {
	vm           *ch8.VirtualMachine
	audioPlayer  *audio.Player
	audioContext *audio.Context
}

// NewSound creates a new instance of the emulator's audio.
func NewSound(vm *ch8.VirtualMachine) *Sound {
	audioContext := audio.NewContext(sampleRate)
	audioPlayer, _ := audio.NewPlayer(audioContext, &stream{})
	audioPlayer.SetVolume(0.25)

	return &Sound{
		vm:           vm,
		audioPlayer:  audioPlayer,
		audioContext: audioContext,
	}
}

// Start starts the sound player.
func (s *Sound) Start() {
	for range time.Tick(16 * time.Millisecond) {
		if s.vm.Sound > 0 {
			s.audioPlayer.Play()
		} else {
			s.audioPlayer.Pause()
		}
	}
}
