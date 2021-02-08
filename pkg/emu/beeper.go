package emu

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// beepStream is a sine wave byte stream played by the CHIP-8 beeper.
//
// This struct is taken directly from Ebiten's example code:
// <https://ebiten.org/examples/sinewave.html>
type beepStream struct {
	frequency  int
	sampleRate int
	position   int64
	remaining  []byte
}

// Read fills the byte stream with sine wave samples.
func (s *beepStream) Read(buf []byte) (int, error) {
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

	length := int64(s.sampleRate / s.frequency)
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

// Close closes the bye stream.
func (s *beepStream) Close() error {
	return nil
}

// Beeper is the CHIP-8 emulator's beeper (audio).
type Beeper struct {
	audioPlayer  *audio.Player
	audioContext *audio.Context
}

// NewBeeper creates a new instance of the emulator's beeper.
func NewBeeper(frequency, sampleRate int) *Beeper {
	audioContext := audio.NewContext(int(sampleRate))
	audioPlayer, _ := audio.NewPlayer(
		audioContext,
		&beepStream{
			frequency:  frequency,
			sampleRate: sampleRate,
		},
	)
	audioPlayer.SetVolume(0.25)

	return &Beeper{
		audioPlayer:  audioPlayer,
		audioContext: audioContext,
	}
}

// Play starts the beeper.
func (s *Beeper) Play() {
	s.audioPlayer.Play()
}

// Stop stops the beeper.
func (s *Beeper) Stop() {
	s.audioPlayer.Pause()
}
