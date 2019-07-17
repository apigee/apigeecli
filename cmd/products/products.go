package products

import (
	"../shared"
	"./crtprod"
	"./delprod"
	"./getprod"
	"./listproducts"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "products",
	Aliases: []string{"prods"},
	Short:   "Manage Apigee API products",
	Long:    "Manage Apigee API products",
}

func init() {

	Cmd.PersistentFlags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listproducts.Cmd)
	Cmd.AddCommand(getprod.Cmd)
	Cmd.AddCommand(delprod.Cmd)
	Cmd.AddCommand(crtprod.Cmd)
}
