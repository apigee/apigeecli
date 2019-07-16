package apps

import (
	"../shared"
	getapp "./getapp"
	listapp "./listapp"
	delapp "./delapp"
	crtapp "./crtapp"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "apps",
	Aliases: []string{"applications"},
	Short: "Manage Apigee Developer Applications",
	Long:  "Manage Apigee Developer Applications",
}

var expand = false
var count string

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listapp.Cmd)
	Cmd.AddCommand(getapp.Cmd)
	Cmd.AddCommand(delapp.Cmd)
	Cmd.AddCommand(crtapp.Cmd)
}
