package targetservers

import (
	"../shared"
	"./crtts"
	"./delts"
	"./getts"
	"./listts"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "targetservers",
	Aliases: []string{"ts"},
	Short:   "Manage Target Servers",
	Long:    "Manage Target Servers",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listts.Cmd)
	Cmd.AddCommand(getts.Cmd)
	Cmd.AddCommand(delts.Cmd)
	Cmd.AddCommand(crtts.Cmd)
}
