package getfh

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a flowhook",
	Long:  "Get a flowhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "flowhooks", name)
		return shared.HttpClient(u.String())
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Flowhook point")

	_ = Cmd.MarkFlagRequired("name")
}
