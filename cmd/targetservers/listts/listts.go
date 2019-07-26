package listts

import (
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List Target Servers",
	Long:  "List Target Servers",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "targetservers")
		return shared.HttpClient(u.String())
	},
}

func init() {

}
