package crtapis

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an API proxy in an Apigee Org",
	Long:  "Creates an API proxy in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")

		if proxy != "" {
			q := u.Query()
			q.Set("name", name)
			q.Set("action", "import")
			u.RawQuery = q.Encode()
			err := shared.ReadBundle(name)
			if err != nil {
				return err
			}

			return shared.PostHttpOctet(u.String(), proxy)
		} else {
			proxyName := "{\"name\":\"" + name + "\"}"
			return shared.HttpClient(u.String(), proxyName)
		}
	},
}

var name, proxy string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	Cmd.Flags().StringVarP(&proxy, "proxy", "p",
		"", "API Proxy Bundle path")

	_ = Cmd.MarkFlagRequired("name")
}
