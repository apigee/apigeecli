package delka

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Key Alias",
	Long:  "Delete a Key Alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores", shared.RootArgs.AliasName, "aliases", aliasName)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
	},
}

var aliasName string

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")

	_ = Cmd.MarkFlagRequired("alias")
}
