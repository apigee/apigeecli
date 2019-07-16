package deldev

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes an App Developer from an organization",
	Long:  "Deletes an App Developer from an organization",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", name)
		shared.HttpClient(u.String(),"","DELETE")
	},
}

var name string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer")

	Cmd.MarkFlagRequired("name")
}
