package listres

import (
	"fmt"
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

//Cmd to list resources
var Cmd = &cobra.Command{
	Use:   "list",
	Short: "List all resources in your environment",
	Long:  "List all resources in your environment",
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		if !types.IsValidResource(resType) {
			return fmt.Errorf("Invalid resource type")
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, "resourcefiles")

		if resType != "" {
			q := u.Query()
			q.Set("type", resType)
			u.RawQuery = q.Encode()
		}
		_, err = shared.HttpClient(true, u.String())
		return

	},
}

var resType string

func init() {

	Cmd.Flags().StringVarP(&resType, "type", "p",
		"", "Resource type")

}
