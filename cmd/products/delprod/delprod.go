package delprod

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to delete products
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an API product from an organization",
	Long:  "Deletes an API product from an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts", name)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "name of the API Product")

	_ = Cmd.MarkFlagRequired("name")
}
