package main

import (
	"fmt"

	"github.com/kevhlee/chip8/pkg/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
