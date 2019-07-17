package listapis

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List APIs in an Apigee Org",
	Long:  "List APIs in an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		if shared.RootArgs.Env != "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "deployments")
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")
		}

		shared.HttpClient(u.String())
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

}
