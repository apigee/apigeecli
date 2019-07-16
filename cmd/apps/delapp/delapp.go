package deldev

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a Developer App from an organization",
	Long:  "Deletes a Developer Appfrom an organization",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", id, "apps", name)
		shared.HttpClient(u.String(),"","DELETE")
	},
}

var name, id string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&id, "id", "i",
		"", "Developer Id")

	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("id")
}
