package delts

import (
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Target Server",
	Long:  "Delete a Target Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "targetservers", name)
		return shared.HttpClient(u.String(), "", "DELETE")
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the targetserver")

	_ = Cmd.MarkFlagRequired("name")
}
