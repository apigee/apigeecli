package listdeploy

import (
	"net/url"
	"../../shared"
	"path"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a shared flow by name",
	Long:  "Gets a shared flow by name, including a list of its revisions.",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows", name)
		shared.GetHttpClient(u.String(), shared.RootArgs.Token)
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow name")

	Cmd.MarkFlagRequired("name")		

}