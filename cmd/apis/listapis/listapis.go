package listapis

import (
	"net/url"
	"path"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to list api
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List APIs in an Apigee Org",
	Long:  "List APIs in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if shared.RootArgs.Env != "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "deployments")
		} else {
			if includeRevisions {
				q := u.Query()
				q.Set("includeRevisions", strconv.FormatBool(includeRevisions))
				u.RawQuery = q.Encode()
			}
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")
		}
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var includeRevisions bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Env, "env", "e",
		"", "Apigee environment name")

	Cmd.Flags().BoolVarP(&includeRevisions, "incl", "i",
		false, "Include revisions")
}
