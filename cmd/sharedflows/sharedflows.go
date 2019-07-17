package sharedflows

import (
	"../shared"
	"./delsf"
	"./depsf"
	"./fetchsf"
	getsf "./getsf"
	listsf "./listsf"
	"./undepsf"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "sharedflows",
	Short: "Manage Apigee shared flows in an org",
	Long:  "Manage Apigee shared flows in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listsf.Cmd)
	Cmd.AddCommand(getsf.Cmd)
	Cmd.AddCommand(fetchsf.Cmd)
	Cmd.AddCommand(delsf.Cmd)
	Cmd.AddCommand(undepsf.Cmd)
	Cmd.AddCommand(depsf.Cmd)
}
