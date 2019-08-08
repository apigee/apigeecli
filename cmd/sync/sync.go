package sync

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/sync/getsync"
	"github.com/srinandan/apigeecli/cmd/sync/setsync"
)

//Cmd to manage identities
var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Manage identities with grant access to control plane resources",
	Long:  "Manage identities with grant access to control plane resources",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(setsync.Cmd)
	Cmd.AddCommand(getsync.Cmd)
}
