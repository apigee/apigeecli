package csr

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to get key aliases
var Cmd = &cobra.Command{
	Use:   "csr",
	Short: "Generates a csr for the private key in an alias",
	Long:  "Generates a PKCS #10 Certificate Signing Request for the private key in an alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env,
			"keystores", shared.RootArgs.AliasName, "aliases", aliasName, "csr")
		_, err = shared.HttpClient(true, u.String())
		return
	},
}

var aliasName string

func init() {

	Cmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("alias")
}
