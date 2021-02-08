package cmd

import (
	"fmt"
	"os"

	"github.com/kevhlee/chip8/pkg/emu"
	"github.com/spf13/cobra"
)

// NewRunCmd creates a new command for running CHIP-8 ROMs.
func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run <path to ROM>",
		Short: "Run a CHIP-8 ROM file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Input a path to a CHIP-8 ROM file")
			}
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}
			return nil
		},
		RunE: run,
	}

	runCmd.Flags().IntP(
		"scale",
		"s",
		emu.DefaultScale,
		"set the scale factor of the CHIP-8 screen",
	)
	runCmd.Flags().Bool(
		"mute",
		false,
		"turn off the sound of the CHIP-8 emulator",
	)
	runCmd.Flags().Bool(
		"debug",
		false,
		"set CHIP-8 emulator to debug mode",
	)

	return runCmd
}

func run(cmd *cobra.Command, args []string) error {
	scale, err := cmd.Flags().GetInt("scale")
	if err != nil {
		return err
	}

	mute, err := cmd.Flags().GetBool("mute")
	if err != nil {
		return err
	}

	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		return err
	}

	emu := emu.NewEmulator(debug, scale, mute)
	if err := emu.LoadROM(args[0]); err != nil {
		return err
	}
	return emu.Start()
}
