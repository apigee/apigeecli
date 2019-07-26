package listenv

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List environments in an Apigee Org",
	Long:  "List environments in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments")
		return shared.HttpClient(u.String())
	},
}

func init() {

}
