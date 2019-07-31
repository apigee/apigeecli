package crtapp

import (
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Developer App",
	Long:  "Create a Developer App",
	RunE: func(cmd *cobra.Command, args []string) error {
		u, _ := url.Parse(shared.BaseURL)

		app := []string{}

		app = append(app, "\"name\":\""+name+"\"")

		if len(apiProducts) > 0 {
			app = append(app, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")
		}

		if callback != "" {
			app = append(app, "\"callbackUrl\":\""+callback+"\"")
		}

		if expires != "" {
			app = append(app, "\"keyExpiresIn\":\""+expires+"\"")
		}

		if len(scopes) > 0 {
			app = append(app, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
		}

		if len(attrs) != 0 {
			attributes := []string{}
			for key, value := range attrs {
				attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
			}
			attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
			app = append(app, attributesStr)
		}

		payload := "{" + strings.Join(app, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", email, "apps")
		return shared.HttpClient(u.String(), payload)
	},
}

var name, email, expires, callback string
var apiProducts, scopes []string
var attrs map[string]string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&email, "email", "e",
		"", "The developer's email or id")
	Cmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in milliseconds, for the lifetime of the consumer key")
	Cmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	Cmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	Cmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	Cmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("email")
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}
