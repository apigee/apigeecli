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
	"encoding/json"
	"io/ioutil"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/products"
	"github.com/apigee/apigeecli/clilog"
	"github.com/spf13/cobra"
)

//Cmd to update a product
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an API product",
	Long:  "Update an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var operationGrpBytes, gqlOperationGrpBytes []byte
		var attributes []products.Attribute

		operationGrp := products.OperationGroup{}
		gqlOperationGrp := products.GraphqlOperationGroup{}

		p := products.Product{}

		p.Name = name
		p.DisplayName = displayName
		p.ApprovalType = approval
		p.Description = description
		p.Quota = quota
		p.QuotaInterval = quotaInterval
		p.QuotaTimeUnit = quotaUnit
		p.Environments = environments
		p.Proxies = proxies
		p.Scopes = scopes

		if operationGroupFile != "" {
			if operationGrpBytes, err = ioutil.ReadFile(operationGroupFile); err != nil {
				clilog.Info.Println(err)
				return err
			}
			if err = json.Unmarshal(operationGrpBytes, &operationGrp); err != nil {
				clilog.Info.Println(err)
				return err
			}
			p.OperationGroup = &operationGrp
		}

		if gqlOperationGroupFile != "" {
			if gqlOperationGrpBytes, err = ioutil.ReadFile(gqlOperationGroupFile); err != nil {
				clilog.Info.Println(err)
				return err
			}
			if err = json.Unmarshal(gqlOperationGrpBytes, &gqlOperationGrp); err != nil {
				clilog.Info.Println(err)
				return err
			}
			p.GraphQLOperationGroup = &gqlOperationGrp
		}

		if len(attrs) > 0 {
			for k, v := range attrs {
				a := products.Attribute{}
				a.Name = k
				a.Value = v
				attributes = append(attributes, a)
			}
			p.Attributes = attributes
		}

		_, err = products.Update(p)

		return
	},
}

var merge bool

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
	UpdateCmd.Flags().StringVarP(&operationGroupFile, "opgrp", "",
		"", "File containing Operation Group JSON. See samples for how to create the file")
	UpdateCmd.Flags().StringVarP(&gqlOperationGroupFile, "gqlopgrp", "",
		"", "File containing GraphQL Operation Group JSON. See samples for how to create the file")
	//TODO: apiresource -r later

	_ = UpdateCmd.MarkFlagRequired("name")
}
