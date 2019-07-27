package crtdev

import (
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a developer",
	Long:  "Create a developer",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)

		developer := []string{}

		developer = append(developer, "\"email\":\""+email+"\"")
		developer = append(developer, "\"firstName\":\""+firstName+"\"")
		developer = append(developer, "\"lastName\":\""+lastName+"\"")
		developer = append(developer, "\"userName\":\""+userName+"\"")

		if len(attrs) != 0 {
			attributes := []string{}
			for key, value := range attrs {
				attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
			}
			attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
			developer = append(developer, attributesStr)
		}

		payload := "{" + strings.Join(developer, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")
		return shared.HttpClient(u.String(), payload)
	},
}

var email, lastName, firstName, userName string
var attrs map[string]string

func init() {

	Cmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	Cmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	Cmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	Cmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")
	Cmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = Cmd.MarkFlagRequired("email")
	_ = Cmd.MarkFlagRequired("first")
	_ = Cmd.MarkFlagRequired("last")
	_ = Cmd.MarkFlagRequired("user")
}
