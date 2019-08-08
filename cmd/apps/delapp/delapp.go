package deldev

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to delete app
var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a Developer App from an organization",
	Long:  "Deletes a Developer Appfrom an organization",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", id, "apps", name)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
	},
}

var name, id string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&id, "id", "i",
		"", "Developer Id")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("id")
}
