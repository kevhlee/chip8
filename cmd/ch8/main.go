package main

import (
	"fmt"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/kevhlee/chip8/pkg/ch8"
	"github.com/spf13/cobra"
)

func logo() string {
	return heredoc.Doc(`
		 ██████╗██╗  ██╗██╗██████╗        █████╗ 
		██╔════╝██║  ██║██║██╔══██╗      ██╔══██╗
		██║     ███████║██║██████╔╝█████╗╚█████╔╝
		██║     ██╔══██║██║██╔═══╝ ╚════╝██╔══██╗
		╚██████╗██║  ██║██║██║           ╚█████╔╝
		 ╚═════╝╚═╝  ╚═╝╚═╝╚═╝            ╚════╝`,
	)
}

func description() string {
	return heredoc.Docf(
		"\n\n%s\n\nA CHIP-8 emulator written in Go.",
		logo(),
	)
}

func main() {
	rootCmd := &cobra.Command{
		Use:  "ch8",
		Long: description(),
	}

	rootCmd.AddCommand(
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
				emu := ch8.NewEmulator()
				emu.LoadROM(args[0])
				return emu.Start()
			},
		},
	)

	rootCmd.Execute()
}
