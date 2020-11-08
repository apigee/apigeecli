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

package envoy

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/products"
)

//CreateCmd
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Envoy binding; Binds an Envoy service to an API Product",
	Long:  "Create a new Envoy binding; Binds an Envoy service to an API Product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if legacy {
			if attrs == nil {
				attrs = make(map[string]string)
			}
			attrs[envoyAttributeName] = strings.Join(serviceNames, ",")
			_, err = products.CreateLegacy(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, remoteServiceProxy, scopes, attrs)
		} else {
			_, err = products.CreateRemoteServiceOperationGroup(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, serviceNames, scopes, attrs)
		}
		return
	},
}

var remoteServiceProxy = []string{"remote-service"}

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	CreateCmd.Flags().StringVarP(&displayName, "displayname", "m",
		"", "Display Name of the API Product")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the API Product")
	CreateCmd.Flags().StringArrayVarP(&environments, "envs", "e",
		[]string{}, "Environments to enable")
	CreateCmd.Flags().StringArrayVarP(&serviceNames, "remote-svcs", "r",
		[]string{}, "Envoy Service names. Ex: -s service1 -s service2")
	CreateCmd.Flags().StringArrayVarP(&scopes, "scopes", "s",
		[]string{}, "OAuth scopes")
	CreateCmd.Flags().StringVarP(&quota, "quota", "q",
		"", "Quota Amount")
	CreateCmd.Flags().StringVarP(&quotaInterval, "interval", "i",
		"", "Quota Interval")
	CreateCmd.Flags().StringVarP(&quotaUnit, "unit", "u",
		"", "Quota Unit")
	CreateCmd.Flags().StringVarP(&approval, "approval", "f",
		"", "Approval type")
	CreateCmd.Flags().StringToStringVar(&attrs, "attrs",
		nil, "Custom attributes")
	CreateCmd.Flags().BoolVarP(&legacy, "legacy", "l",
		false, "Legacy product object")

	_ = CreateCmd.MarkFlagRequired("prod")
	_ = CreateCmd.MarkFlagRequired("remote-svcs")
	_ = CreateCmd.MarkFlagRequired("approval")
}
