package main

import (
	"fmt"

	"github.com/kevhlee/chip8/pkg/cmd"
)

func main() {
	if err := cmd.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
	}
}
