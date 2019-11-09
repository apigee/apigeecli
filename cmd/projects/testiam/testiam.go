package testiam

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "testiam",
	Short: "Test IAM policy for a GCP Project",
	Long:  "Test IAM policy for a GCP Project",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.CrmURL)
		u.Path = path.Join(u.Path, shared.RootArgs.ProjectID+":testIamPermissions")
		payload := "{\"permissions\":[\"resourcemanager.projects.get\"]}"
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.ProjectID, "prj", "p",
		"", "GCP Project ID")

	_ = Cmd.MarkFlagRequired("prj")
}
