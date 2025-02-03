package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kevhlee/chip8/chip8"
)

func main() {
	var opts chip8.Options

	flag.IntVar(&opts.Scale, "scale", 8, "set the scaling factor of the screen")
	flag.IntVar(&opts.TPS, "tps", 12, "set the number of CPU ticks per frame")
	flag.Parse()

	filename := flag.Arg(0)
	if len(filename) == 0 {
		exit("Usage: chip8 <path to ROM>")
	}

	ch8 := chip8.New(opts)
	if err := ch8.LoadROM(flag.Arg(0)); err != nil {
		exit(err.Error())
	}

	if err := ch8.Run(); err != nil {
		exit(err.Error())
	}
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
