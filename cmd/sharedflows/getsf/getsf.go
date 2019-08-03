package getsf

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a shared flow by name",
	Long:  "Gets a shared flow by name, including a list of its revisions.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows", name)
		_, err = shared.HttpClient(true, u.String()) 
		return 
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow name")

	_ = Cmd.MarkFlagRequired("name")

}
