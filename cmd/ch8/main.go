package main

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/kevhlee/chip8/pkg/emu"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:     "ch8",
		Example: heredoc.Doc(`$ ch8 run roms/Logo.ch8`),
		Long: heredoc.Docf(`

			 ██████╗██╗  ██╗██╗██████╗        █████╗ 
			██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
			██║     ███████║██║██████╔╝█████╗╚█████╔╝
			██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
			╚██████╗██║  ██║██║██║           ╚█████╔╝
			 ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝

			A CHIP-8 emulator written in Go.
		`),
	}

	cmd.AddCommand(
		&cobra.Command{
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
				emu.LoadROM(args[0])
				return emu.Start()
			},
		},
	)

	cmd.Execute()
}
