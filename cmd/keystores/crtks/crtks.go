package crtks

import (
	"github.com/spf13/cobra"
)

//Cmd to create key stores
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Store",
	Long:  "Create a Key Store",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}
