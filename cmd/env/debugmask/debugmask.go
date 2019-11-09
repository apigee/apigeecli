package debugmask

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/env/debugmask/getdebugmask"
	"github.com/srinandan/apigeecli/cmd/env/debugmask/setdebugmask"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "debugmask",
	Short: "Manage debugmasks for the environment",
	Long:  "Manage debugmasks for the environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(getdebugmask.Cmd)
	Cmd.AddCommand(setdebugmask.Cmd)
}
