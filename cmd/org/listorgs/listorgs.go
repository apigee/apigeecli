package listorgs

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List the Apigee organizations",
	Long:  "List the Apigee organizations, and the related projects that a user has permissions for",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		shared.GetHttpClient(u.String(), shared.RootArgs.Token)
	},
}

func init() {

}
