package getcert

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get key aliases
var Cmd = &cobra.Command{
	Use:   "getcert",
	Short: "Get a Key alias certificate",
	Long:  "Get a Key alias certificate",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"keystores", shared.RootArgs.AliasName, "aliases", aliasName, "certificate")
		err = shared.DownloadResource(u.String(), aliasName+".crt", "")
		return
	},
}

var aliasName string

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("alias")
}
