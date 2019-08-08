package getorg

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

//Cmd to get org details
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Show details of an Apigee Org",
	Long:  "Show details of an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkFlagRequired("org")
}
