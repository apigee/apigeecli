package gettk

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new access token",
	Long:  "Generate a new access token",
	Args: func(cmd *cobra.Command, args []string) error {
		if shared.RootArgs.ServiceAccount == "" {
			return fmt.Errorf("service account cannot be empty")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		shared.Init()
		err := shared.SetAccessToken()
		if err != nil {
			return err
		}
		fmt.Println(shared.RootArgs.Token)
		return nil
	},
}

func init() {

}
