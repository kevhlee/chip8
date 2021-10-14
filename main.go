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

	command.Flags().IntVarP(
		&options.Scale,
		"scale",
		"s",
		ch8.DefaultScale,
		"set the scale factor of the CHIP-8 screen",
	)

	command.Flags().Float64VarP(
		&options.Volume,
		"volume",
		"v",
		ch8.DefaultVolume,
		"set the volume of the CHIP-8 emulator",
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err)
	}
}
