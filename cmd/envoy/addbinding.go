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

//AddBindCmd
var AddBindCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new Envoy binding; Binds an Envoy service to an existing API Product",
	Long:  "Add a new Envoy binding; Binds an Envoy service to an existing API Product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		apiclient.SetPrintOutput(false)
		_, err = products.GetAttribute(productName, envoyAttributeName)
		apiclient.SetPrintOutput(true)
		if err != nil {
			attr := make(map[string]string)
			attr[string(envoyAttributeName)] = strings.Join(serviceNames, ",")
			_, err = products.UpdateLegacy(productName, "", "", "", "", "", "", nil, nil, nil, attr)
			return err
		} else {
			_, err = products.UpdateAttribute(productName, envoyAttributeName, strings.Join(serviceNames, ","))
			return err
		}

		return err
	},
}

func init() {
	AddBindCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	AddBindCmd.Flags().StringVarP(&productName, "prod", "p",
		"", "Apigee API Product name")
	AddBindCmd.Flags().StringArrayVarP(&serviceNames, "remote-svcs", "r",
		[]string{}, "Envoy Service names")

	_ = AddBindCmd.MarkFlagRequired("prod")
	_ = AddBindCmd.MarkFlagRequired("remote-svcs")
}
