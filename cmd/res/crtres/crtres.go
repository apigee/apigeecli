package crtres

import (
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

//Cmd to create a resource
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource file",
	Long:  "Create a resource file",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {

		if !types.IsValidResource(resType) {
			return fmt.Errorf("invalid resource type")
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "resourcefiles")
		if resType != "" {
			q := u.Query()
			q.Set("type", resType)
			q.Set("name", name)
			u.RawQuery = q.Encode()
		}
		_, err = shared.PostHttpOctet(true, u.String(), resPath)
		return
	},
}

var name, resType, resPath string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the resource file")
	Cmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type")
	Cmd.Flags().StringVarP(&resPath, "respath", "r",
		"", "Resource Path")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("type")
	_ = Cmd.MarkFlagRequired("respath")
}
