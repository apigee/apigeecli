package targetservers

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/flowhooks/crtfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/delfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/getfh"
	"github.com/srinandan/apigeecli/cmd/flowhooks/listfh"
	"github.com/srinandan/apigeecli/cmd/shared"
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
