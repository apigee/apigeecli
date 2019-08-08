package listks

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list key stores
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List Key Stores",
	Long:  "List Key Stores",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

}
