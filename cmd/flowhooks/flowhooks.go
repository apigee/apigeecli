package targetservers

import (
	"../shared"
	"./crtfh"
	"./delfh"
	"./getfh"
	"./listfh"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "flowhooks",
	Aliases: []string{"fs"},
	Short:   "Manage Flowhooks",
	Long:    "Manage Flowhooks",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listfh.Cmd)
	Cmd.AddCommand(getfh.Cmd)
	Cmd.AddCommand(delfh.Cmd)
	Cmd.AddCommand(crtfh.Cmd)
}
