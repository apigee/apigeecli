package apps

import (
	"../shared"
	crtapp "./crtapp"
	delapp "./delapp"
	genkey "./genkey"
	getapp "./getapp"
	listapp "./listapp"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"applications"},
	Short:   "Manage Apigee Developer Applications",
	Long:    "Manage Apigee Developer Applications",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listapp.Cmd)
	Cmd.AddCommand(getapp.Cmd)
	Cmd.AddCommand(delapp.Cmd)
	Cmd.AddCommand(crtapp.Cmd)
	Cmd.AddCommand(genkey.Cmd)
}
