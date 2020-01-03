// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package crtprod

import (
	"net/url"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to create a new product
var Cmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API product",
	Long:  "Create an API product",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
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
		if len(attrs) != 0 {
			attributes := []string{}
			for key, value := range attrs {
				attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
			}
			attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
			product = append(product, attributesStr)
		}

		payload := "{" + strings.Join(product, ",") + "}"
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apiproducts")
		_, err = shared.HttpClient(true, u.String(), payload)
		return
	},
}

var name, description, approval, displayName, quota, quotaInterval, quotaUnit string
var environments, proxies, scopes []string
var attrs map[string]string

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
	Cmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	//TODO: apiresource -r later

	_ = Cmd.MarkFlagRequired("name")
	_ = Cmd.MarkFlagRequired("envs")
	_ = Cmd.MarkFlagRequired("proxies")
	_ = Cmd.MarkFlagRequired("approval")
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}
