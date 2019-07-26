package listfh

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List Flowhooks",
	Long:  "List Flowhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "flowhooks")
		return shared.HttpClient(u.String())
	},
}

func init() {

}
