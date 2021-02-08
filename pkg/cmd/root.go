package cmd

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

// NewRootCmd creates a new root command for running the CHIP-8
// emulator.
func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
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

	rootCmd.AddCommand(NewRunCmd())

	return rootCmd
}
