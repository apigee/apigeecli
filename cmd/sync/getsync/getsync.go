package getsync

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Show the list of identities with access to control plane resources",
	Long:  "Show the list of identities with access to control plane resources",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org+":getSyncAuthorization")
		shared.PostHttpClient(u.String(), "")
	},
}

func init() {

}
