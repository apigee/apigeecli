package listsf

import (
	"net/url"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list shared flow
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all shared flows in the organization.",
	Long:  "Lists all shared flows in the organization.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if shared.RootArgs.Env != "" {
			q := u.Query()
			q.Set("sharedFlows", "true")
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "deployments")
		} else {
			if revisions {
				q := u.Query()
				q.Set("includeRevisions", strconv.FormatBool(revisions))
				u.RawQuery = q.Encode()
			}
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows")
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var revisions bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")
	Cmd.Flags().BoolVarP(&revisions, "rev", "r",
		false, "Include revisions")

}
