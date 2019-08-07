package getks

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Key Store",
	Long:  "Get a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores", name)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("name")
}
