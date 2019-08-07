package main

import (
	"os"

	"github.com/srinandan/apigeecli/cmd"
)

var Version string

func main() {

	rootCmd := cmd.GetRootCmd()
	rootCmd.Version = "0.5, Git: " + Version

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
