package setmart

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
)

var Cmd = &cobra.Command{
	Use:   "setmart",
	Short: "Set MART endpoint for an Apigee Org",
	Long:  "Set MART endpoint for an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		orgname := "\"name\":\"" + shared.RootArgs.Org + "\","
		martprop := "{\"name\":\"features.mart.server.endpoint\","
		martpropvalue := "\"value\":\"" + mart + "\"}"
		props := "\"properties\": {" + "\"property\": [" + martprop + martpropvalue + "]}"
		payload := "{" + orgname + props + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org)
		_ = shared.HttpClient(u.String(), payload)
	},
}

var mart string

func init() {

	Cmd.Flags().StringVarP(&shared.RootArgs.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.Flags().StringVarP(&mart, "mart", "m",
		"", "MART Endpoint")

	_ = Cmd.MarkFlagRequired("org")
}
