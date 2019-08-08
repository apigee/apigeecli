package listenv

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list envs
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List environments in an Apigee Org",
	Long:  "List environments in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

}
