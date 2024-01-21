package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kevhlee/chip8/ch8"
)

func main() {
	opts := ch8.EmulatorOptions{}

	flag.IntVar(&opts.Scale, "scale", ch8.DefaultScale, "set the scale factor of the screen")
	flag.Parse()

	romPath := flag.Arg(0)
	if len(romPath) == 0 {
		exitWithError("Usage: chip8 <path to ROM>")
	}

	emu, err := ch8.NewEmulator(opts)
	if err != nil {
		exitWithError(err.Error())
	}

	if err := emu.LoadROM(romPath); err != nil {
		exitWithError(err.Error())
	}

	if err := emu.Start(); err != nil {
		if err != ch8.ErrTerminated {
			exitWithError(err.Error())
		}
	}
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
