package delkvm

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a KVM map",
	Long:  "Delete a KVM map",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keyvaluemaps", name)

		_ = shared.HttpClient(u.String(), "", "DELETE")
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Environment name")
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "KVM Map name")

	_ = Cmd.MarkFlagRequired("env")
	_ = Cmd.MarkFlagRequired("name")
}
