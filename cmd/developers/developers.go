package developers

import (
	"../shared"
	"./crtdev"
	"./deldev"
	"./getdev"
	"./listdev"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "devs",
	Aliases: []string{"developers"},
	Short:   "Manage Apigee App Developers",
	Long:    "Manage Apigee App Developers",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listdev.Cmd)
	Cmd.AddCommand(getdev.Cmd)
	Cmd.AddCommand(deldev.Cmd)
	Cmd.AddCommand(crtdev.Cmd)
}
