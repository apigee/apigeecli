package listdeploy

import (
	"fmt"
	"net/url"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list deployed api
var Cmd = &cobra.Command{
	Use:   "listdeploy",
	Short: "Lists all deployments of an API proxy",
	Long:  "Lists all deployments of an API proxy",
	Args: func(cmd *cobra.Command, args []string) error {
		if shared.RootArgs.Env != "" && revision == -1 {
			return fmt.Errorf("proxy revision must be supplied with the environment")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if shared.RootArgs.Env != "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "apis", name, "revisions", strconv.Itoa(revision), "deployments")
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis", name, "deployments")
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string
var revision int

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().IntVarP(&revision, "rev", "r",
		-1, "API Proxy revision")

	_ = Cmd.MarkFlagRequired("name")
}
