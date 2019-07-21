package main

import (
	"./cmd"
	"os"
)

var GitVersion string

func main() {

	rootCmd := cmd.GetRootCmd()
	rootCmd.Version = "0.2, Git: " + GitVersion

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
