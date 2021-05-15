package main

import (
	"fmt"
	"micro/cmd"
	"os"
)

// root execute command with cobra
func main() {
	if err := cmd.Runner.RootCmd().Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run command: %v\n", err)
		os.Exit(1)
	}
}