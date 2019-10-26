package apis

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/apis/traceapi/crttrcapi"
	"github.com/srinandan/apigeecli/cmd/apis/traceapi/gettrcapi"
	"github.com/srinandan/apigeecli/cmd/apis/traceapi/listtrcapi"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "trace",
	Short: "Manage debugging/tracing of Apigee API proxies",
	Long:  "Manage debugging/tracing of Apigee API proxy revisions deployed in an environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(crttrcapi.Cmd)
	Cmd.AddCommand(listtrcapi.Cmd)
	Cmd.AddCommand(gettrcapi.Cmd)
}
