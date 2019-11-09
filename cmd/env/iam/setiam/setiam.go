package setiam

import (
	"encoding/json"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "set",
	Short: "Set an IAM policy on an Environment",
	Long:  "Set an IAM policy on an Environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//the following steps will validate json
		m := map[string]string{}
		err = json.Unmarshal([]byte(payload), &m)
		if err != nil {
			return err
		}
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env+":setIamPolicy")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var payload string

func init() {

	Cmd.Flags().StringVarP(&payload, "policy", "p",
		"", "IAM policy is in JSON format")

	_ = Cmd.MarkFlagRequired("policy")
}
