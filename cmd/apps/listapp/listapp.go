package listapp

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of Developer Applications",
	Long:  "Returns a list of app IDs within an organization based on app status",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apps")
		q := u.Query()
		if expand {
			q.Set("expand", "true")
		} else {
			q.Set("expand", "false")
		}
		if count != "" {
			q.Set("row", count)
		}
		u.RawQuery = q.Encode()
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var expand = false
var count string

func init() {

	Cmd.Flags().StringVarP(&count, "count", "c",
		"", "Number of apps; limit is 1000")

	Cmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")
}
