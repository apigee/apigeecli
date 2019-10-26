package crtka

import (
	"github.com/spf13/cobra"
)

//Cmd to create key aliases
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Key Alias",
	Long:  "Create a Key Alias",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
}
