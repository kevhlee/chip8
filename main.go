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

	interpreter := chip8.NewInterpreter(opts)
	if err := interpreter.LoadROM(flag.Arg(0)); err != nil {
		exit(err.Error())
	}

	if err := interpreter.Run(); err != nil {
		exit(err.Error())
	}
}

func exit(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}
