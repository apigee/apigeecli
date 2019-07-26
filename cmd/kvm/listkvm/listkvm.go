package listkvm

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of KVMs",
	Long:  "Returns a list of KVMs",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keyvaluemaps")

		_ = shared.HttpClient(u.String())
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Environment name")
	_ = Cmd.MarkFlagRequired("env")
}
