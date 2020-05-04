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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/products"
)

//RemoveCmd
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes an Envoy binding from an API Product",
	Long:  "Removes an Envoy binding from an API Product",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		fmt.Println("Current Values of Attribute: ")
		respBody, err := products.GetAttribute(productName, envoyAttributeName)
		if err != nil {
			return err
		}

		var attributeValue map[string]string
		err = json.Unmarshal(respBody, &attributeValue)
		if err != nil {
			return err
		}

		values := strings.Split(attributeValue["value"], ",")
		var newValues []string
		var found bool

		for _, value := range values {
			if strings.TrimSpace(value) != serviceName {
				newValues = append(newValues, strings.TrimSpace(value))
			} else {
				found = true
				fmt.Println(("Found service name, removing binding"))
			}
		}

		if !found {
			return fmt.Errorf("did not find value matching service name")
		}

		fmt.Println("New Values of Attribute: ")
		newAttributeValues := strings.Join(newValues, ",")
		_, err = products.UpdateAttribute(productName, envoyAttributeName, newAttributeValues)
		return err
	},
}

var serviceName string

func init() {

	RemoveCmd.Flags().StringVarP(&productName, "prod", "p",
		"", "Apigee API Product name")
	RemoveCmd.Flags().StringVarP(&serviceName, "svc", "s",
		"", "Envoy Service name")

	_ = RemoveCmd.MarkFlagRequired("prod")
	_ = RemoveCmd.MarkFlagRequired("svc")
}
