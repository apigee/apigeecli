package keyaliases

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/keyaliases/crtka"
	"github.com/srinandan/apigeecli/cmd/keyaliases/delka"
	"github.com/srinandan/apigeecli/cmd/keyaliases/getka"
	"github.com/srinandan/apigeecli/cmd/keyaliases/listka"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage key aliases
var Cmd = &cobra.Command{
	Use:     "keyaliases",
	Aliases: []string{"ka"},
	Short:   "Manage Key Aliases",
	Long:    "Manage Key Aliases",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.AliasName, "name", "n",
		"", "Key Store name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	_ = Cmd.MarkPersistentFlagRequired("env")
	_ = Cmd.MarkPersistentFlagRequired("name")

	Cmd.AddCommand(listka.Cmd)
	Cmd.AddCommand(getka.Cmd)
	Cmd.AddCommand(delka.Cmd)
	Cmd.AddCommand(crtka.Cmd)
}
