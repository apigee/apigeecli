package crtapp

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create a app App",
	Long:  "Create a app App",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		app := []string{}

		app = append(app, "\"name\":\""+name+"\"")
		app = append(app, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")

		if callback != "" {
			app = append(app, "\"callbackUrl\":\""+callback+"\"")
		}

		if expires != "" {
			app = append(app, "\"keyExpiresIn\":\""+expires+"\"")		
		}

		if len(scopes) > 0 {
			app = append(app, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
		}
		
		payload := "{"+strings.Join(app,",")+"}"	
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", email, "apps")
		shared.HttpClient(u.String(), payload)
	},
}

var name, email, expires, callback string
var apiProducts, scopes []string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the app app")	
	Cmd.Flags().StringVarP(&email, "email", "e",
		"", "The app's email")
	Cmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in milliseconds, for the lifetime of the consumer key")		
	Cmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")		
	Cmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	Cmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")

	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("email")
	Cmd.MarkFlagRequired("prods")
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str,",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}