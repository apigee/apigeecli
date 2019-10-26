package res

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/res/listres"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage resources
var Cmd = &cobra.Command{
	Use:     "resources",
	Aliases: []string{"res"},
	Short:   "Manage resources within an Apigee environment",
	Long:    "Manage resources within an Apigee environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listres.Cmd)
}
