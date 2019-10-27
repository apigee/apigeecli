package crttrcapi

import (
	"bytes"
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new debug session for an API proxy",
	Long:  "Create a new debug session for Apigee API proxy revision deployed in an environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"apis", name, "revisions", revision, "debugsessions")
		q := u.Query()
		q.Set("timeout", "567")
		u.RawQuery = q.Encode()

		var payload string
		if len(filter) != 0 {
			payload = "{\"filter\":" + getFilterStr() + "}"
		} else {
			payload = "{}"
		}
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var name, revision string
var filter map[string]string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API proxy name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")
	Cmd.Flags().StringToStringVar(&filter, "filter",
		nil, "Filter Conditions; format is name1=value1,name2=value2...")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
}

func getFilterStr() string {
	b := new(bytes.Buffer)
	for key, value := range filter {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}
