package testiam

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "test",
	Short: "Test IAM policy for an Environment",
	Long:  "Test IAM policy for an Environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env+":testIamPermissions")
		payload := "{\"permissions\":[\"apigee.environments.get\"]}"
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}
