package fetchapi

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to download api
var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted proxy bundle ",
	Long:  "Returns a zip-formatted proxy bundle of code and config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return shared.FetchBundle("apis", name, revision)
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy Bundle Name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "API Proxy revision")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
}
