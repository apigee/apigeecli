package getdev

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Returns the profile for a developer by email address or ID",
	Long:  "Returns the profile for a developer by email address or ID",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "email of the developer")

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("name")
}
