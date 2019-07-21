package gettk

import (
	"../../shared"
	"fmt"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate a new access token",
	Long:  "Generate a new access token",
	Args: func(cmd *cobra.Command, args []string) error {
		if shared.RootArgs.ServiceAccount != "" {
			return nil
		} else {
			return fmt.Errorf("Service account cannot be empty")
		}
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
