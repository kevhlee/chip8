package chip8

type Event int

const (
	EventIgnore Event = iota
	EventQuit
	EventPause
	EventReset
	EventNextCycle
)
