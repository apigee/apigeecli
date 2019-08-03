package getapi

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets an API Proxy by name",
	Long:  "Gets an API Proxy by name, including a list of its revisions.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis", name)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")

	_ = Cmd.MarkFlagRequired("name")

}
