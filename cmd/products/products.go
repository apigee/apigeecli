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
	"os"

	"internal/clilog"

	"internal/client/products"

	"github.com/spf13/cobra"
)

// Cmd to manage products
var Cmd = &cobra.Command{
	Use:     "products",
	Aliases: []string{"prods"},
	Short:   "Manage Apigee API products",
	Long:    "Manage Apigee API products and Rate Plans for Monetization",
}

var org, name string
var conn int

var description, approval, displayName, quota, quotaInterval, quotaUnit string
var environments, proxies, scopes []string
var attrs map[string]string

func init() {

	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DelCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(ImpCmd)
	Cmd.AddCommand(ExpCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(AttributesCmd)
	Cmd.AddCommand(RatePlanCmd)
}

func getOperationGroup(operationGroupFile string) (*products.OperationGroup, error) {
	var operationGrpBytes []byte
	var err error

	operationGrp := products.OperationGroup{}

	if operationGroupFile != "" {
		if operationGrpBytes, err = os.ReadFile(operationGroupFile); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if err = json.Unmarshal(operationGrpBytes, &operationGrp); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		return &operationGrp, nil
	}
	return nil, nil
}

func getGqlOperationGroup(gqlOperationGroupFile string) (*products.GraphqlOperationGroup, error) {
	var gqlOperationGrpBytes []byte
	var err error

	gqlOperationGrp := products.GraphqlOperationGroup{}

	if gqlOperationGroupFile != "" {
		if gqlOperationGrpBytes, err = os.ReadFile(gqlOperationGroupFile); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if err = json.Unmarshal(gqlOperationGrpBytes, &gqlOperationGrp); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		return &gqlOperationGrp, nil
	}
	return nil, nil
}

func getGrpcOperationGroup(grpcOperationGroupFile string) (*products.GrpcOperationGroup, error) {
	var grpcOperationGrpBytes []byte
	var err error

	grpcOperationGrp := products.GrpcOperationGroup{}

	if gqlOperationGroupFile != "" {
		if grpcOperationGrpBytes, err = os.ReadFile(gqlOperationGroupFile); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		if err = json.Unmarshal(grpcOperationGrpBytes, &grpcOperationGrp); err != nil {
			clilog.Info.Println(err)
			return nil, err
		}
		return &grpcOperationGrp, nil
	}
	return nil, nil
}

func getAttributes(attrs map[string]string) []products.Attribute {
	var attributes []products.Attribute

	if len(attrs) > 0 {
		for k, v := range attrs {
			a := products.Attribute{}
			a.Name = k
			a.Value = v
			attributes = append(attributes, a)
		}
		return attributes
	}

	return nil
}
