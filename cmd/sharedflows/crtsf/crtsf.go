package crtsf

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a sharedflow in an Apigee Org",
	Long:  "Creates a sharedflow in an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows")

		if proxy != "" {
			q := u.Query()
			q.Set("name", name)
			q.Set("action", "import")
			u.RawQuery = q.Encode()
			err = shared.ReadBundle(name)
			if err != nil {
				return err
			}
			_, err = shared.PostHttpOctet(true, u.String(), proxy)
			return
		}
		proxyName := "{\"name\":\"" + name + "\"}"
		_, err = shared.HttpClient(true, u.String(), proxyName)
		return
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
