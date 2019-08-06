package apis

import (
	"github.com/spf13/cobra"
	crtapi "github.com/srinandan/apigeecli/cmd/apis/crtapi"
	delapi "github.com/srinandan/apigeecli/cmd/apis/delapi"
	"github.com/srinandan/apigeecli/cmd/apis/depapi"
	"github.com/srinandan/apigeecli/cmd/apis/expapis"
	fetch "github.com/srinandan/apigeecli/cmd/apis/fetchapi"
	impapis "github.com/srinandan/apigeecli/cmd/apis/impapis"
	"github.com/srinandan/apigeecli/cmd/apis/listapis"
	"github.com/srinandan/apigeecli/cmd/apis/listdeploy"
	"github.com/srinandan/apigeecli/cmd/apis/undepapi"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API proxies in an org",
	Long:  "Manage Apigee API proxies in an org",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listapis.Cmd)
	Cmd.AddCommand(listdeploy.Cmd)
	Cmd.AddCommand(crtapi.Cmd)
	Cmd.AddCommand(expapis.Cmd)
	Cmd.AddCommand(depapi.Cmd)
	Cmd.AddCommand(delapi.Cmd)
	Cmd.AddCommand(fetch.Cmd)
	Cmd.AddCommand(impapis.Cmd)
	Cmd.AddCommand(undepapi.Cmd)
}
