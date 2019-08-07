package keystores

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/keystores/crtks"
	"github.com/srinandan/apigeecli/cmd/keystores/delks"
	"github.com/srinandan/apigeecli/cmd/keystores/getks"
	"github.com/srinandan/apigeecli/cmd/keystores/listks"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:     "keystores",
	Aliases: []string{"ks"},
	Short:   "Manage Key Stores",
	Long:    "Manage Key Stores",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(listks.Cmd)
	Cmd.AddCommand(getks.Cmd)
	Cmd.AddCommand(delks.Cmd)
	Cmd.AddCommand(crtks.Cmd)
}
