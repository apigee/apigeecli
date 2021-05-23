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
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/products"
	"github.com/srinandan/apigeecli/clilog"
)

//Cmd to create a new product
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API product",
	Long:  "Create an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if legacy && (operationGroupFile != "" || gqlOperationGroupFile != "") {
			return fmt.Errorf("operationGroupFile/gqlOperationGroupFile cannot be used with legacy mode")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if legacy {
			_, err = products.CreateLegacy(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, proxies, scopes, attrs)
		} else {
			var operationGrp, gqlOperationGrp []byte
			if operationGroupFile != "" {
				operationGrp, err = ioutil.ReadFile(operationGroupFile)
				if err != nil {
					clilog.Info.Println(err)
					return err
				}
			}
			if gqlOperationGroupFile != "" {
				gqlOperationGrp, err = ioutil.ReadFile(gqlOperationGroupFile)
				if err != nil {
					clilog.Info.Println(err)
					return err
				}
			}
			_, err = products.CreateProxyOperationGroup(name, description, approval, displayName, quota, quotaInterval, quotaUnit, environments, scopes, operationGrp, gqlOperationGrp, attrs)
		}
		return
	},
}

var legacy bool
var operationGroupFile, gqlOperationGroupFile string

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	CreateCmd.Flags().StringVarP(&displayName, "displayname", "m",
		"", "Display Name of the API Product")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the API Product")
	CreateCmd.Flags().StringArrayVarP(&environments, "envs", "e",
		[]string{}, "Environments to enable")
	CreateCmd.Flags().StringArrayVarP(&proxies, "proxies", "p",
		[]string{}, "API Proxies in product. Ex: -p proxy1 -p proxy2")
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
	CreateCmd.Flags().StringVarP(&operationGroupFile, "opgrp", "",
		"", "File containing Operation Group JSON. See samples for how to create the file")
	CreateCmd.Flags().StringVarP(&gqlOperationGroupFile, "gqlopgrp", "",
		"", "File containing GraphQL Operation Group JSON. See samples for how to create the file")
	CreateCmd.Flags().BoolVarP(&legacy, "legacy", "l",
		false, "Legacy product object")

	//TODO: apiresource -r later

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("approval")
}
