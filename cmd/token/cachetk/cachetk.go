package cachetk

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to cache token
var Cmd = &cobra.Command{
	Use:   "cache",
	Short: "Generate and cache a new access token",
	Long:  "Generate and cache a new access token",
	Args: func(cmd *cobra.Command, args []string) error {
		if shared.RootArgs.ServiceAccount == "" {
			return fmt.Errorf("Service account cannot be empty")
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		shared.Init()
		token, err := shared.GenerateAccessToken()
		fmt.Printf("Token %s cached\n", token)
		return err
	},
}

func init() {

}
