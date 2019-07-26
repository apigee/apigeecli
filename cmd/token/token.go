package token

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/token/gettk"
)

var Cmd = &cobra.Command{
	Use:   "token",
	Short: "Manage OAuth 2.0 access tokens",
	Long:  "Manage OAuth 2.0 access tokens",
}

func init() {

	Cmd.AddCommand(gettk.Cmd)
}
