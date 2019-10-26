package delcache

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to delete cache
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a cache resource from the environment",
	Long:  "Delete a cache resource from the environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "caches", name)
		_, err = shared.HttpClient(true, u.String())
		return

	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")

	_ = Cmd.MarkFlagRequired("name")

}
