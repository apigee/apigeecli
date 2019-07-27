package getdev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get App in an Organization by App ID",
	Long:  "Returns the app profile for the specified app ID.",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps", appId)
		return shared.HttpClient(u.String())
	},
}

var appId string

func init() {

	Cmd.Flags().StringVarP(&appId, "appId", "i",
		"", "Developer app id")

	_ = Cmd.MarkFlagRequired("appId")
}
