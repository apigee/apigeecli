package createsync

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get org details
var Cmd = &cobra.Command{
	Use:   "createsync",
	Short: "Create a new IAM Service Account for Apigee Synchronizer",
	Long:  "Create a new IAM Service Account for Apigee Synchronizer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return shared.CreateIAMServiceAccount(name, "sync")
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.ProjectID, "prj", "p",
		"", "GCP Project ID")
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Service Account Name")

	_ = Cmd.MarkFlagRequired("prj")
	_ = Cmd.MarkFlagRequired("name")
}
