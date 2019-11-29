package enableac

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/org/setprop"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to set mart endpoint
var Cmd = &cobra.Command{
	Use:   "enable-apigee-connect",
	Short: "Set MART endpoint for an Apigee Org",
	Long:  "Set MART endpoint for an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return setprop.SetOrgProperty("features.mart.apigee.connect.enabled", "true")
	},
}

var mart string
var whitelist bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")

	_ = Cmd.MarkFlagRequired("org")
}
