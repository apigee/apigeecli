package listenv

import (
	"net/url"
	"../../shared"
	"path"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List environments in an Apigee Org",
	Long:  "List environments in an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments")
		shared.GetHttpClient(u.String(), shared.RootArgs.Token)
	},
}

func init() {

}