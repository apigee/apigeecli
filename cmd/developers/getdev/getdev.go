package getdev

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Returns the profile for a developer by email address or ID",
	Long:  "Returns the profile for a developer by email address or ID",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name)
		_ = shared.HttpClient(u.String())
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
