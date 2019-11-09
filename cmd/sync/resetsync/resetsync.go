package resetsync

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to set identities
var Cmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset identities with access to control plane resources",
	Long:  "Reset identities with access to control plane resources",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org+":setSyncAuthorization")
		payload := "{\"identities\":[]}"
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}
