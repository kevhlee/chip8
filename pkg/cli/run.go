package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/kevhlee/chip8/pkg/ch8"
	"github.com/spf13/cobra"
)

// NewRunCmd creates a new command for running CHIP-8 ROMs.
func NewRunCmd() *cobra.Command {
	runCmd := &cobra.Command{
		Use:   "run <path to ROM>",
		Short: "Run a CHIP-8 ROM file",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("input a path to a CHIP-8 ROM file")
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
		ch8.DefaultScale,
		"set the scale factor of the CHIP-8 screen",
	)
	runCmd.Flags().Float64P(
		"volume",
		"v",
		0.25,
		"set the volume of the CHIP-8 emulator",
	)

	return runCmd
}

func run(cmd *cobra.Command, args []string) error {
	scale, err := cmd.Flags().GetInt("scale")
	if err != nil {
		return err
	} else if scale < 1 {
		return errors.New("scale factor must be positive")
	}

	volume, err := cmd.Flags().GetFloat64("volume")
	if err != nil {
		return err
	} else if volume < 0.0 || volume > 1.0 {
		return errors.New("volume must be between [0.0, 1.0]")
	}

	emu := ch8.NewEmulator(scale, volume)
	if err := emu.LoadROM(args[0]); err != nil {
		return err
	}
	return emu.Start()
}
