package apis

import (
	"../shared"
	crtapi "./crtapi"
	delapi "./delapi"
	"./depapi"
	fetch "./fetchapi"
	"./listapis"
	"./listdeploy"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API proxies in an org",
	Long:  "Manage Apigee API proxies in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listapis.Cmd)
	Cmd.AddCommand(listdeploy.Cmd)
	Cmd.AddCommand(crtapi.Cmd)
	Cmd.AddCommand(depapi.Cmd)
	Cmd.AddCommand(delapi.Cmd)
	Cmd.AddCommand(fetch.Cmd)
}
