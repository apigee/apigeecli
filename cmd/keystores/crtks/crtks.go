package crtks

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create key stores
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Store",
	Long:  "Create a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "keystores")
		payload := "{\"name\":\"" + name + "\"}"
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the key store")

	_ = Cmd.MarkFlagRequired("name")
}