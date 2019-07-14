package apis

import (
	"./listapis"
	"../shared"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API proxies deployed to an environment",
	Long:  "Manage Apigee API proxies deployed to an environment",
}

var expand = false
var count string

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
	"", "Apigee environment name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.MarkPersistentFlagRequired("env")
	Cmd.AddCommand(listapis.Cmd)		
}