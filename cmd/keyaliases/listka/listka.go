package listka

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List Key Aliases",
	Long:  "List Key Alises",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores", shared.RootArgs.AliasName, "aliases")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

}
