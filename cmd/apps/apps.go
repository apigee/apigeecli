package apps

import (
	listapp "./listapp"
	getapp "./getapp"
	"../shared"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "apps",
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
}