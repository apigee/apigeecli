package kvm

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/kvm/crtkvm"
	"github.com/srinandan/apigeecli/cmd/kvm/delkvm"
	"github.com/srinandan/apigeecli/cmd/kvm/listkvm"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "kvms",
	Short: "Manage Key Value Maps",
	Long:  "Manage Key Value Maps",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listkvm.Cmd)
	Cmd.AddCommand(delkvm.Cmd)
	Cmd.AddCommand(crtkvm.Cmd)
}
