package delkvm

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

//Cmd to delete kvm
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a KVM map",
	Long:  "Delete a KVM map",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keyvaluemaps", name)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
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
