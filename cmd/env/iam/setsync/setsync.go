package setsync

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "setsync",
	Short: "Set Synchronization Manager role for a SA on an environment",
	Long:  "Set Synchronization Manager role for a SA on an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return shared.SetIAMServiceAccount(serviceAccountName, "sync")
	},
}

var serviceAccountName string

func init() {

	Cmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")

	_ = Cmd.MarkFlagRequired("name")
}
