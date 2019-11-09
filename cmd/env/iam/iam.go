package iam

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/env/iam/getiam"
	"github.com/srinandan/apigeecli/cmd/env/iam/setiam"
	"github.com/srinandan/apigeecli/cmd/env/iam/testiam"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "iam",
	Short: "Manage IAM permissions for the environment",
	Long:  "Manage IAM permissions for the environment",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	_ = Cmd.MarkPersistentFlagRequired("env")

	Cmd.AddCommand(getiam.Cmd)
	Cmd.AddCommand(setiam.Cmd)
	Cmd.AddCommand(testiam.Cmd)
}
