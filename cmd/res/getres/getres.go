package getres

import (
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

//Cmd to get a resource
var Cmd = &cobra.Command{
	Use:   "get",
	Short: "Get a resource file",
	Long:  "Get a resource file",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if !types.IsValidResource(resType) {
			return fmt.Errorf("invalid resource type")
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "resourcefiles", resType, name)
		err = shared.DownloadResource(u.String(), name, resType)
		return
	},
}

var name, resType string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the resource file")
	Cmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("type")
}
