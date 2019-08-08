package env

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/env/getenv"
	"github.com/srinandan/apigeecli/cmd/env/listenv"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage envs
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
