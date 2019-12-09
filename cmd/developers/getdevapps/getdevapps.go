package getdevapps

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

//Cmd to get developer
var Cmd = &cobra.Command{
	Use:   "getapps",
	Short: "Returns the apps owned by a developer by email address",
	Long:  "Returns the apps owned by a developer by email address",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if expand {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name, "apps")
			q := u.Query()
			q.Set("expand", "true")
			u.RawQuery = q.Encode()
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name, "apps")
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var name string
var expand bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "email of the developer")

	Cmd.Flags().BoolVarP(&expand, "expand", "x",
		false, "expand app details")

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("name")
}
