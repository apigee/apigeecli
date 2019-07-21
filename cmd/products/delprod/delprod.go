package delprod

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an API product from an organization",
	Long:  "Deletes an API product from an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts", name)
		return shared.HttpClient(u.String(), "", "DELETE")
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "name of the API Product")

	_ = Cmd.MarkFlagRequired("name")
}
