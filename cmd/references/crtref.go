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

package references

import (
	"fmt"

	"internal/apiclient"

	"github.com/apigee/apigeecli/client/references"
	"github.com/spf13/cobra"
)

// CreateCmd to create references
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a reference in an environment",
	Long:  "Create a reference in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if resourceType != "KeyStore" && resourceType != "TrustStore" {
			return fmt.Errorf("resourceType should be KeyStore or TrustStore")
		}
		apiclient.SetApigeeEnv(env)

		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = references.Create(name, description, resourceType, refers)
		return
	},
}

func init() {

	CreateCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "Reference name")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the reference")
	CreateCmd.Flags().StringVarP(&resourceType, "restype", "s",
		"", "Resource type must be KeyStore or TrustStore")
	CreateCmd.Flags().StringVarP(&refers, "refers", "r",
		"", "The id of the resource to which this reference refers")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("restype")
	_ = CreateCmd.MarkFlagRequired("refers")
}
