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

	"internal/client/references"

	"github.com/spf13/cobra"
)

// UpdateCmd to update references
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a reference in an environment",
	Long:  "Update a reference in an environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if resourceType != "KeyStore" && resourceType != "TrustStore" {
			return fmt.Errorf("resourceType should be KeyStore or TrustStore")
		}
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = references.Update(name, description, resourceType, refers)
		return
	},
}

func init() {
	UpdateCmd.PersistentFlags().StringVarP(&name, "name", "n",
		"", "Reference name")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the reference")
	UpdateCmd.Flags().StringVarP(&resourceType, "restype", "s",
		"", "Resource type must be KeyStore or TrustStore")
	UpdateCmd.Flags().StringVarP(&refers, "refers", "r",
		"", "The id of the resource to which this reference refers")

	_ = UpdateCmd.MarkFlagRequired("name")
}
