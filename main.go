package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/kevhlee/chip8/ch8"
	"github.com/spf13/cobra"
)

func main() {
	cli := &cobra.Command{
		Use:     "ch8",
		Example: heredoc.Doc(`$ ch8 run roms/Logo.ch8`),
		Long: heredoc.Doc(
			`

			 ██████╗██╗  ██╗██╗██████╗        █████╗ 
			██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
			██║     ███████║██║██████╔╝█████╗╚█████╔╝
			██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
			╚██████╗██║  ██║██║██║           ╚█████╔╝
			 ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝

			A CHIP-8 emulator written in Go.
			`,
		),
	}

	cli.AddCommand(newRunCmd())

	if err := cli.Execute(); err != nil {
		fmt.Println(err)
	}
}

func newRunCmd() *cobra.Command {
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

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			scale, _ := cmd.Flags().GetInt("scale")
			volume, _ := cmd.Flags().GetFloat64("volume")

			emu := ch8.NewEmulator(scale, volume)

			if err := emu.LoadROM(args[0]); err != nil {
				return err
			}

			return emu.Start()
		},
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
		ch8.DefaultVolume,
		"set the volume of the CHIP-8 emulator",
	)

	return runCmd
}
