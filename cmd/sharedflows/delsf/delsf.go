package delsf

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a shared flow",
	Long:  "Deletes a shared flow and all associated policies, resources, and revisions. The flow must be undeployed first.",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "sharedflows", name)
		_, err = shared.HttpClient(true, u.String(), "", "DELETE")
		return
	},
}

var name string

func init() {
	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "shared flow name")

	_ = Cmd.MarkFlagRequired("name")

}
