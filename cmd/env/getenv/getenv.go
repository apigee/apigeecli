package getenv

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get properties of an environment",
	Long:  "Get properties of an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)

		if config {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "deployedConfig")
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env)
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var config = false

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().BoolVarP(&config, "config", "c",
		false, "Return configuration details")

	_ = Cmd.MarkFlagRequired("env")
}
