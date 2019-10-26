package gettrcapi

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a debug session for an API proxy revision",
	Long:  "Get a debug session for an API proxy revision deployed in an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		if messageID == "" {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "apis", name, "revisions", revision, "debugsessions", sessionID, "data")
			q := u.Query()
			q.Set("limit", "20")
			u.RawQuery = q.Encode()
		} else {
			u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "apis", name, "revisions", revision, "debugsessions", sessionID, "data", messageID)
		}

		_, err = shared.HttpClient(true, u.String())
		return

	},
}

var name, revision, sessionID, messageID string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")
	Cmd.Flags().StringVarP(&sessionID, "ses", "s",
		"", "Debug session Id")
	Cmd.Flags().StringVarP(&messageID, "msg", "m",
		"", "Debug session Id")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
	_ = Cmd.MarkFlagRequired("ses")

}
