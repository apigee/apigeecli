package crtprod

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"strings"
)

var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API product",
	Long:  "Create an API product",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)

		product := []string{}

		product = append(product, "\"name\":\""+name+"\"")
		
		if displayName == "" {
			product = append(product, "\"displayName\":\""+name+"\"")
		} else {
			product = append(product, "\"displayName\":\""+displayName+"\"")
		}
		
		if description != "" {
			product = append(product, "\"description\":\""+description+"\"")
		}
		product = append(product, "\"environments\":[\""+getArrayStr(environments)+"\"]")
		product = append(product, "\"proxies\":[\""+getArrayStr(proxies)+"\"]")

		if len(scopes) > 0 {
			product = append(product, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
		}

		product = append(product, "\"approvalType\":\""+approval+"\"")

		if quota != "" {
			product = append(product, "\"quota\":\""+quota+"\"")
		}
		if quotaInterval != "" {
			product = append(product, "\"quotaInterval\":\""+quotaInterval+"\"")
		}
		if quotaUnit != "" {
			product = append(product, "\"quotaTimeUnit\":\""+quotaUnit+"\"")
		}
		payload := "{"+strings.Join(product,",")+"}"	
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts")
		shared.HttpClient(u.String(), payload)
	},
}

var name, description, approval, displayName, quota, quotaInterval, quotaUnit string
var environments, proxies, scopes []string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	Cmd.Flags().StringVarP(&displayName, "displayname", "m",
		"", "Display Name of the API Product")		
	Cmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the API Product")		
	Cmd.Flags().StringArrayVarP(&environments, "envs", "e",
		[]string{}, "Environments to enable")
	Cmd.Flags().StringArrayVarP(&proxies, "proxies", "p",
		[]string{}, "API Proxies in product")
	Cmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	Cmd.Flags().StringVarP(&quota, "quota", "q",
		"", "Quota Amount")
	Cmd.Flags().StringVarP(&quotaInterval, "interval", "i",
		"", "Quota Interval")
	Cmd.Flags().StringVarP(&quotaUnit, "unit", "u",
		"", "Quota Unit")		
	Cmd.Flags().StringVarP(&approval, "approval", "f",
		"", "Approval type")
	//TODO: apiresource -r later	

	Cmd.MarkFlagRequired("name")
	Cmd.MarkFlagRequired("envs")
	Cmd.MarkFlagRequired("proxies")
	Cmd.MarkFlagRequired("approval")
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str,",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}