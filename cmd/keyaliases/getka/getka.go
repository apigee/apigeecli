package getka

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get key aliases
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Key Store",
	Long:  "Get a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores", shared.RootArgs.AliasName, "aliases", aliasName)
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var aliasName string

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("alias")
}
