package crtsf

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a sharedflow in an Apigee Org",
	Long:  "Creates a sharedflow in an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows")

		if proxy != "" {
			q := u.Query()
			q.Set("name", name)
			q.Set("action", "import")
			u.RawQuery = q.Encode()
			err := shared.ReadBundle(name)
			if err == nil {
				_ = shared.PostHttpOctet(u.String(), proxy)
			}
		} else {
			proxyName := "{\"name\":\"" + name + "\"}"
			_ = shared.HttpClient(u.String(), proxyName)
		}
	},
}

var name, proxy string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	Cmd.Flags().StringVarP(&proxy, "proxy", "p",
		"", "Sharedflow bundle path")

	_ = Cmd.MarkFlagRequired("name")
}
