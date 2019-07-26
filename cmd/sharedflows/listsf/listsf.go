package listsf

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all shared flows in the organization.",
	Long:  "Lists all shared flows in the organization.",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		if shared.RootArgs.Env != "" {
			q := u.Query()
			q.Set("sharedFlows", "true")
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "deployments")
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows")
		}
		return shared.HttpClient(u.String())
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

}
