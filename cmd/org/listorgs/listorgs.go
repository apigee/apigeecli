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
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		return shared.HttpClient(u.String())
	},
}

func init() {

}
