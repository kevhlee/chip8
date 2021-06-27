package main

import (
	"fmt"

	"github.com/kevhlee/chip8/pkg/cli"
)

func main() {
	rootCmd := cli.NewRootCmd()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
