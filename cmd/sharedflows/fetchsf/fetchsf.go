package fetchsf

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to download shared flow
var Cmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted shared flow bundle ",
	Long:  "Returns a zip-formatted shared flow bundle of code and config files",
	RunE: func(cmd *cobra.Command, args []string) error {
		return shared.FetchBundle("sharedflows", name, revision)
	},
}

var name, revision string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow Bundle Name")
	Cmd.Flags().StringVarP(&revision, "rev", "v",
		"", "Shared flow revision")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("rev")
}
