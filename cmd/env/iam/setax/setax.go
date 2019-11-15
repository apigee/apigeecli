package setax

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to manage tracing of apis
var Cmd = &cobra.Command{
	Use:   "setax",
	Short: "Set Analytics Agent role for a SA on an environment",
	Long:  "Set Analytics Agent role for a SA an Environment",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return shared.SetIAMServiceAccount(serviceAccountName, "analytics")
	},
}

var serviceAccountName string

func init() {

	Cmd.Flags().StringVarP(&serviceAccountName, "name", "n",
		"", "Service Account Name")

	_ = Cmd.MarkFlagRequired("name")
}
