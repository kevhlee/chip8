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
		Short: "Run a CHIP-8 ROM",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("Input a path to a CHIP-8 ROM file")
			}
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			emu := emu.NewEmulator()
			if err := emu.LoadROM(args[0]); err != nil {
				return err
			}
			return emu.Start()
		},
	}

	return runCmd
}
