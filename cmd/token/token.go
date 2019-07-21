package token

import (
	"./gettk"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "token",
	Short: "Manage OAuth 2.0 access tokens",
	Long:  "Manage OAuth 2.0 access tokens",
}

func init() {

	Cmd.AddCommand(gettk.Cmd)
}
