package setmart

import (
	"net/url"
	"path"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "setmart",
	Short: "Set MART endpoint for an Apigee Org",
	Long:  "Set MART endpoint for an Apigee Org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		orgname := "\"name\":\"" + shared.RootArgs.Org + "\","
		hybridprop := "{\"name\" : \"features.hybrid.enabled\", \"value\" : \"true\"},"
		martprop := "{\"name\":\"features.mart.server.endpoint\", \"value\":\"" + mart + "\"}"

		props := "\"properties\": {" + "\"property\": [" + hybridprop + martprop

		if whitelist {
			whitelistprop := "{\"name\":\"features.mart.server.endpoint\", \"value\" : \"true\"}"
			props = props + "," + whitelistprop
		}

		props = props + "]}"

		payload := "{" + orgname + props + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org)
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var mart string
var whitelist bool

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.Flags().StringVarP(&mart, "mart", "m",
		"", "MART Endpoint")
	Cmd.Flags().BoolVarP(&whitelist, "whitelist", "w",
		false, "Whitelist GCP Source IP")

	_ = Cmd.MarkFlagRequired("org")
	_ = Cmd.MarkFlagRequired("mart")
}
