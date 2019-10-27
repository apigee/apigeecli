package main

import (
	"os"

	"github.com/srinandan/apigeecli/cmd"
)

//Version contains the git hash
var Version string

func main() {

	rootCmd := cmd.GetRootCmd()
	rootCmd.Version = "0.7, Git: " + Version

	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
