package getsync

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

//Cmd to get list of identities
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Show the list of identities with access to control plane resources",
	Long:  "Show the list of identities with access to control plane resources",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org+":getSyncAuthorization")
		_, err = shared.HttpClient(true, u.String(), "")
		return
	},
}

func init() {

}
