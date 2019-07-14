package main

import (
	"os"

	"./cmd"
)

func main() {
	rootCmd := cmd.GetRootCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}