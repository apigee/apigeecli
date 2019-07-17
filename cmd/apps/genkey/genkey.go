package crtkey

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a new developer KeyPair",
	Long:  "Generate a new developer KeyPair",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		key := []string{}

		key = append(key, "\"name\":\""+name+"\"")
		key = append(key, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")

		if callback != "" {
			key = append(key, "\"callbackUrl\":\""+callback+"\"")
		}

		if expires != "" {
			key = append(key, "\"keyExpiresIn\":\""+expires+"\"")
		}

		if len(scopes) > 0 {
			key = append(key, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
		}

		payload := "{" + strings.Join(key, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "developers", devid, "apps", name)
		shared.HttpClient(u.String(), payload)
	},
}

var name, devid, expires, callback string
var apiProducts, scopes []string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the developer app")
	Cmd.Flags().StringVarP(&devid, "devid", "d",
		"", "Developer Id")
	Cmd.Flags().StringVarP(&expires, "expires", "x",
		"", "A setting, in milliseconds, for the lifetime of the consumer key")
	Cmd.Flags().StringVarP(&callback, "callback", "c",
		"", "The callbackUrl is used by OAuth")
	Cmd.Flags().StringArrayVarP(&apiProducts, "prods", "p",
		[]string{}, "A list of api products")
	Cmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")

	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("devid")
	Cmd.MarkFlagRequired("prods")
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}
