package delks

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to delete key stores
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Key Store",
	Long:  "Delete a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores", name)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("name")
}
