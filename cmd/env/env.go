package env

import (
	"../shared"
	"./getenv"
	"./listenv"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "envs",
	Aliases: []string{"environments"},
	Short:   "Manage Apigee environments",
	Long:    "Manage Apigee environments",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listenv.Cmd)
	Cmd.AddCommand(getenv.Cmd)
}
