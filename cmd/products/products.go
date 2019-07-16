package products

import (
	"../shared"
	"./getprod"
	"./listproducts"
	"./delprod"
	"./crtprod"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "prods",
	Aliases: []string{"products"},
	Short: "Manage Apigee API products",
	Long:  "Manage Apigee API products",
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
