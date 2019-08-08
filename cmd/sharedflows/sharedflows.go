package sharedflows

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/sharedflows/crtsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/delsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/depsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/expsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/fetchsf"
	getsf "github.com/srinandan/apigeecli/cmd/sharedflows/getsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/impsfs"
	listsf "github.com/srinandan/apigeecli/cmd/sharedflows/listsf"
	"github.com/srinandan/apigeecli/cmd/sharedflows/undepsf"
)

//Cmd to manage shared flows
var Cmd = &cobra.Command{
	Use:   "sharedflows",
	Short: "Manage Apigee shared flows in an org",
	Long:  "Manage Apigee shared flows in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listsf.Cmd)
	Cmd.AddCommand(getsf.Cmd)
	Cmd.AddCommand(fetchsf.Cmd)
	Cmd.AddCommand(delsf.Cmd)
	Cmd.AddCommand(undepsf.Cmd)
	Cmd.AddCommand(depsf.Cmd)
	Cmd.AddCommand(crtsf.Cmd)
	Cmd.AddCommand(expsf.Cmd)
	Cmd.AddCommand(impsfs.Cmd)
}
