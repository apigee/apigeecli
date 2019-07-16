package crtdev

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a developer",
	Long:  "Create a developer",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		developer := []string{}

		developer = append(developer, "\"email\":\""+email+"\"")
		developer = append(developer, "\"firstName\":\""+firstName+"\"")
		developer = append(developer, "\"lastName\":\""+lastName+"\"")
		developer = append(developer, "\"userName\":\""+userName+"\"")

		payload := "{" + strings.Join(developer, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers")
		shared.HttpClient(u.String(), payload)
	},
}

var email, lastName, firstName, userName string

func init() {

	Cmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	Cmd.Flags().StringVarP(&firstName, "first", "f",
		"", "The first name of the developer")
	Cmd.Flags().StringVarP(&lastName, "last", "s",
		"", "The last name of the developer")
	Cmd.Flags().StringVarP(&userName, "user", "u",
		"", "The username of the developer")

	Cmd.MarkFlagRequired("email")
	Cmd.MarkFlagRequired("first")
	Cmd.MarkFlagRequired("last")
	Cmd.MarkFlagRequired("user")
}
