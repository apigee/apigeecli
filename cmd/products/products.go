package products

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/products/crtprod"
	"github.com/srinandan/apigeecli/cmd/products/delprod"
	"github.com/srinandan/apigeecli/cmd/products/getprod"
	"github.com/srinandan/apigeecli/cmd/products/impprod"
	"github.com/srinandan/apigeecli/cmd/products/listproducts"
	"github.com/srinandan/apigeecli/cmd/shared"
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

	_ = Cmd.MarkPersistentFlagRequired("org")

	Cmd.AddCommand(listproducts.Cmd)
	Cmd.AddCommand(getprod.Cmd)
	Cmd.AddCommand(delprod.Cmd)
	Cmd.AddCommand(crtprod.Cmd)
	Cmd.AddCommand(impprod.Cmd)
}
