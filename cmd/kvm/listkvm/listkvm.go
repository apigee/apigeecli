package listkvm

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list kvms
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of KVMs",
	Long:  "Returns a list of KVMs",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keyvaluemaps")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Environment name")
	_ = Cmd.MarkFlagRequired("env")
}
