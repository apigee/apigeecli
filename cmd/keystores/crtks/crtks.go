package crtks

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Store",
	Long:  "Create a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}

func init() {

}
