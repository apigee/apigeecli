package apps

import (
	"github.com/spf13/cobra"
	crtapp "github.com/srinandan/apigeecli/cmd/apps/crtapp"
	delapp "github.com/srinandan/apigeecli/cmd/apps/delapp"
	expapp "github.com/srinandan/apigeecli/cmd/apps/expapp"
	genkey "github.com/srinandan/apigeecli/cmd/apps/genkey"
	getapp "github.com/srinandan/apigeecli/cmd/apps/getapp"
	impapps "github.com/srinandan/apigeecli/cmd/apps/impapps"
	listapp "github.com/srinandan/apigeecli/cmd/apps/listapp"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage apps
var Cmd = &cobra.Command{
	Use:     "apps",
	Aliases: []string{"applications"},
	Short:   "Manage Apigee Developer Applications",
	Long:    "Manage Apigee Developer Applications",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listapp.Cmd)
	Cmd.AddCommand(getapp.Cmd)
	Cmd.AddCommand(delapp.Cmd)
	Cmd.AddCommand(crtapp.Cmd)
	Cmd.AddCommand(genkey.Cmd)
	Cmd.AddCommand(impapps.Cmd)
	Cmd.AddCommand(expapp.Cmd)
}
