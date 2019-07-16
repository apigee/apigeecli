package delapis

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an API proxy",
	Long:  "Deletes an API proxy and all associated endpoints, policies, resources, and revisions. The proxy must be undeployed first.",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis", name)
		shared.HttpClient(u.String(), "", "DELETE")
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")

	Cmd.MarkFlagRequired("name")

}
