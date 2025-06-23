// Copyright 2025 Google LLC
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

package attributes

import (
	"internal/apiclient"
	"internal/client/hub"
	"os"

	"github.com/spf13/cobra"
)

// UpdateCmd to update a catalog items
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an Attribute",
	Long:  "Update an Attribute",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		var aValues []byte

		if aValuesPath != "" {
			if aValues, err = os.ReadFile(aValuesPath); err != nil {
				return err
			}
		}
		_, err = hub.UpdateAttribute(attributeID, updateMask, displayName, description, scope, dataType, aValues, cardinality)
		return
	},
}

var updateMask string

func init() {
	UpdateCmd.Flags().StringVarP(&attributeID, "id", "i",
		"", "Attribute ID")
	UpdateCmd.Flags().StringVarP(&updateMask, "mask", "m",
		"", "Update Mask")
	UpdateCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Attribute Display Name")
	UpdateCmd.Flags().StringVarP(&description, "description", "",
		"", "Attribute Description")
	UpdateCmd.Flags().StringVarP(&scope, "scope", "s",
		"", "Attribute scope")
	UpdateCmd.Flags().StringVarP(&dataType, "data-type", "",
		"", "Attribute data type")
	UpdateCmd.Flags().IntVarP(&cardinality, "cardinality", "c",
		1, "Attribute cardinality")
	UpdateCmd.Flags().StringVarP(&aValuesPath, "allowed-values", "",
		"", "Path to a file containing allowed values")

	_ = UpdateCmd.MarkFlagRequired("id")
	_ = UpdateCmd.MarkFlagRequired("display-name")
	_ = UpdateCmd.MarkFlagRequired("scope")
	_ = UpdateCmd.MarkFlagRequired("data-type")
	_ = UpdateCmd.MarkFlagRequired("mask")
}
