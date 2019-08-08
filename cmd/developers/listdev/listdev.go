package listdev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list developer
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Returns a list of App Developers",
	Long:  "Lists all developers in an organization by email address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")
		q := u.Query()
		if expand {
			q.Set("expand", "true")
		} else {
			q.Set("expand", "false")
		}
		if count != "" {
			q.Set("count", count)
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
		"", "Number of developers; limit is 1000")

	Cmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "Expand Details")
}
