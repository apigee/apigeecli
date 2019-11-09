package projects

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/projects/testiam"
)

//Cmd to manage orgs
var Cmd = &cobra.Command{
	Use:   "projects",
	Short: "Manage GCP Projects that have Apigee enabled",
	Long:  "Manage GCP Projects that have Apigee enabled",
}

func init() {

	Cmd.AddCommand(testiam.Cmd)
}
