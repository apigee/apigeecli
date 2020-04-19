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

package products

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/client/products"
)

//Cmd to update a product
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an API product",
	Long:  "Update an API product",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = products.Update(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, attrs)
		return
	},
}

func init() {

	UpdateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	UpdateCmd.Flags().StringVarP(&displayName, "displayname", "m",
		"", "Display Name of the API Product")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the API Product")
	UpdateCmd.Flags().StringArrayVarP(&environments, "envs", "e",
		[]string{}, "Environments to enable")
	UpdateCmd.Flags().StringArrayVarP(&proxies, "proxies", "p",
		[]string{}, "API Proxies in product")
	UpdateCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	UpdateCmd.Flags().StringVarP(&quota, "quota", "q",
		"", "Quota Amount")
	UpdateCmd.Flags().StringVarP(&quotaInterval, "interval", "i",
		"", "Quota Interval")
	UpdateCmd.Flags().StringVarP(&quotaUnit, "unit", "u",
		"", "Quota Unit")
	UpdateCmd.Flags().StringVarP(&approval, "approval", "f",
		"", "Approval type")
	UpdateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")

	//TODO: apiresource -r later

	_ = UpdateCmd.MarkFlagRequired("name")
	_ = UpdateCmd.MarkFlagRequired("envs")
	_ = UpdateCmd.MarkFlagRequired("proxies")
	_ = UpdateCmd.MarkFlagRequired("approval")
}
