package developers

import (
	"../shared"
	"./getdev"
	"./listdev"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "devs",
	Short: "Manage Apigee App Developers",
	Long:  "Manage Apigee App Developers",
}

var expand = false
var count string

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listdev.Cmd)
	Cmd.AddCommand(getdev.Cmd)
}
