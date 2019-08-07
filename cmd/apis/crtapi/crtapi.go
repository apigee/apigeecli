package crtapis

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create api
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an API proxy in an Apigee Org",
	Long:  "Creates an API proxy in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if proxy != "" {
			return shared.ImportBundle("apis", name, proxy)
		}
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")
		proxyName := "{\"name\":\"" + name + "\"}"
		_, err = shared.HttpClient(true, u.String(), proxyName)
		return
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
