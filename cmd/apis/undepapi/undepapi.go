package undepapi

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "undeploy",
	Short: "Undeploys a revision of an existing API proxy",
	Long:  "Undeploys a revision of an existing API proxy to an environment in an organization",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "apis", name, "revisions", revision, "deployments")
		shared.HttpClient(u.String(), "", "DELETE")
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")

	Cmd.MarkFlagRequired("env")
	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("rev")

}
