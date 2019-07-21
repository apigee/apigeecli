package getapi

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an API Proxy by name",
	Long:  "Gets an API Proxy by name, including a list of its revisions.",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis", name)
		_ = shared.HttpClient(u.String())
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")

	_ = Cmd.MarkFlagRequired("name")

}
