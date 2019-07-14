package products

import (
	"./listproducts"
	"../shared"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "products",
	Short: "Manage Apigee API products",
	Long:  "Manage Apigee API products",
}

var expand = false
var count string

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")
	Cmd.AddCommand(listproducts.Cmd)		
}