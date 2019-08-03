package getprod

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an API product from an organization",
	Long:  "Gets an API product from an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts", name)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "name of the API Product")

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("name")
}
