package listorgs

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
)

//Cmd to list orgs
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List the Apigee organizations",
	Long:  "List the Apigee organizations, and the related projects that a user has permissions for",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

}
