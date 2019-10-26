package listcache

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list caches
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List all caches in your environment",
	Long:  "List all caches in your environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "caches")
		_, err = shared.HttpClient(true, u.String())
		return

	},
}

func init() {

}
