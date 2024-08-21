// Copyright 2020-2024 Google LLC
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
	"internal/apiclient"
	"internal/client/products"
	"strings"

	"github.com/spf13/cobra"
)

// CreateCmd to create a new product
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an API product",
	Long:  "Create an API product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		p := products.APIProduct{}

		// if the user provides scope as comma separated values then handle it
		if len(scopes) != 0 && strings.Contains(scopes[0], ",") {
			scopes = append(scopes[1:], strings.Split(scopes[0], ",")...)
		}

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

		p.OperationGroup, err = getOperationGroup(operationGroupFile)
		if err != nil {
			return err
		}

		p.GraphQLOperationGroup, err = getGqlOperationGroup(gqlOperationGroupFile)
		if err != nil {
			return nil
		}

		p.GrpcOperationGroup, err = getGrpcOperationGroup(grpcOperationGroupFile)
		if err != nil {
			return nil
		}

		p.Attributes = getAttributes(attrs)

		if quotaCounterScope != "" {
			p.QuotaCounterScope = quotaCounterScope
		}

		_, err = products.Create(p)

		return err
	},
	Example: `Create an API product, using externally defined operations, with auto approval and access=public:
` + GetExample(0) + `

Create an API product, using externally defined operations, with auto approval and access=public, and a quota of 100/minute:
` + GetExample(1),
}

var operationGroupFile, gqlOperationGroupFile, grpcOperationGroupFile, quotaCounterScope string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the API Product")
	CreateCmd.Flags().StringVarP(&displayName, "display-name", "m",
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
	CreateCmd.Flags().StringVarP(&grpcOperationGroupFile, "grpcopgrp", "",
		"", "File containing gRPC Operation Group JSON. See samples for how to create the file")
	CreateCmd.Flags().StringVarP(&quotaCounterScope, "quota-counter-scope", "",
		"", "Scope of the quota decides how the quota counter gets applied; can be PROXY or OPERATION")
	// TODO: apiresource -r later

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("display-name")
	_ = CreateCmd.MarkFlagRequired("approval")
}
