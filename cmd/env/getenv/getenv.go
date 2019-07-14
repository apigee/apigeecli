package getenv

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get properties of an environment",
	Long:  "Get properties of an environment",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		if config {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments",shared.RootArgs.Env, "deployedConfig")
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments",shared.RootArgs.Env)
		}
		shared.GetHttpClient(u.String(), shared.RootArgs.Token)
	},

}

var config = false

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
	"", "Apigee environment name")
	Cmd.Flags().BoolVarP(&config, "config", "c",
		false, "Return configuration details")

	Cmd.MarkFlagRequired("env")
}