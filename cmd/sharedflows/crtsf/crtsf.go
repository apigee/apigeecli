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
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows")

		if proxy != "" {
			q := u.Query()
			q.Set("name", name)
			q.Set("action", "import")
			u.RawQuery = q.Encode()
			err := shared.ReadBundle(name)
			if err == nil {
				return shared.PostHttpOctet(u.String(), proxy)
			} else {
				return err
			}
		} else {
			proxyName := "{\"name\":\"" + name + "\"}"
			return shared.HttpClient(u.String(), proxyName)
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
