package env

import (
	"./listenv"
	"./getenv"
	"../shared"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "envs",
	Short: "Manage Apigee environments",
	Long:  "Manage Apigee environments",
}

var expand = false
var count string

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listenv.Cmd)
	Cmd.AddCommand(getenv.Cmd)		
}