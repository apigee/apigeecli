package depapi

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a revision of an existing API proxy",
	Long:  "Deploys a revision of an existing API proxy to an environment in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		if overrides {
			q := u.Query()
			q.Set("override", "true")
			u.RawQuery = q.Encode()
		}
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Org, "apis", name, "revisions", revision, "deployments")
		shared.HttpClient(u.String())
	},
}

var name, revision string
var overrides bool

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")
	Cmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")

	Cmd.MarkFlagRequired("env")
	Cmd.MarkFlagRequired("name")

}
