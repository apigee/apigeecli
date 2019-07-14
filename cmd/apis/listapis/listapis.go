package listapis

import (
	"net/url"
	"../../shared"
	"path"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List APIs in an Apigee Org",
	Long:  "List APIs in an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")
		shared.GetHttpClient(u.String(), shared.RootArgs.Token)
	},
}

func init() {

}