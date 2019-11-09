package getiam

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the IAM policy on an Environment",
	Long:  "Gets the IAM policy on an Environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env+":getIamPolicy")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}
