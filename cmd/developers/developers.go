package developers

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/developers/crtdev"
	"github.com/srinandan/apigeecli/cmd/developers/deldev"
	"github.com/srinandan/apigeecli/cmd/developers/expdev"
	"github.com/srinandan/apigeecli/cmd/developers/getdev"
	"github.com/srinandan/apigeecli/cmd/developers/impdev"
	"github.com/srinandan/apigeecli/cmd/developers/listdev"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:     "devs",
	Aliases: []string{"developers"},
	Short:   "Manage Apigee App Developers",
	Long:    "Manage Apigee App Developers",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listdev.Cmd)
	Cmd.AddCommand(getdev.Cmd)
	Cmd.AddCommand(deldev.Cmd)
	Cmd.AddCommand(crtdev.Cmd)
	Cmd.AddCommand(impdev.Cmd)
	Cmd.AddCommand(expdev.Cmd)
}
