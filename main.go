package main

import (
	"fmt"
	"os"

	"github.com/kevhlee/chip8/ch8"
	"github.com/spf13/cobra"
)

func main() {
	options := ch8.NewEmulatorOptions()

	command := &cobra.Command{
		Use:     "ch8",
		Example: "$ ch8 roms/Logo.ch8",
		Long:    "A CHIP-8 emulator written in Go.",
		Args: func(cli *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("input a path to a CHIP-8 ROM file")
			} else if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			emu, err := ch8.NewEmulator(options)
			if err != nil {
				return err
			}

			if err := emu.LoadROM(args[0]); err != nil {
				return err
			}

			return emu.Start()
		},
	}

	initFlags(command, options)

	if err := command.Execute(); err != nil {
		fmt.Println(err)
	}
}

func initFlags(command *cobra.Command, options *ch8.EmulatorOptions) {
	command.Flags().DurationVar(
		&options.HertzIO,
		"hertz-io",
		ch8.DefaultHertzIO,
		"set the speed of the IO timers.",
	)

	command.Flags().DurationVar(
		&options.HertzVM,
		"hertz-vm",
		ch8.DefaultHertzVM,
		"set the speed of the virtual machine's CPU cycle.",
	)

	command.Flags().IntVarP(
		&options.MaxTPS,
		"tps",
		"t",
		ch8.DefaultMaxTPS,
		"set the max ticks-per-second (TPS) of the renderer",
	)

	command.Flags().IntVarP(
		&options.Scale,
		"scale",
		"s",
		ch8.DefaultScale,
		"set the scale factor of the screen",
	)

	command.Flags().Float64VarP(
		&options.Volume,
		"volume",
		"v",
		ch8.DefaultVolume,
		"set the volume of the emulator",
	)
}
