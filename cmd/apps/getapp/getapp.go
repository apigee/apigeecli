package getdev

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps", name)
		_ = shared.HttpClient(u.String())
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of of the developer app")

	_ = Cmd.MarkFlagRequired("name")
}
