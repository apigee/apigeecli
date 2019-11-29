package enablewl

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/org/setprop"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to set mart endpoint
var Cmd = &cobra.Command{
	Use:   "enable-mart-whitelist",
	Short: "Enable IP whitelisting for MART connections",
	Long:  "Enable IP whitelisting for MART connections",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return setprop.SetOrgProperty("features.mart.ip.whitelist.enabled", "true")
	},
}

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkFlagRequired("org")
}
