// Copyright 2024 Google LLC
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
	"os"

	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// CrtCmd
var CrtCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new Attribute in API Hub",
	Long:  "Create a new Attribute in API Hub",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		var aValues []byte

		if aValues, err = os.ReadFile(aValuesPath); err != nil {
			return err
		}
		_, err = hub.CreateAttribute(attributeID, displayName, description, scope, dataType, aValues, cardinality)
		return
	},
}

var (
	attributeID, displayName, description, scope, dataType, aValuesPath string
	cardinality                                                         int
)

func init() {
	CrtCmd.Flags().StringVarP(&attributeID, "id", "i",
		"", "Attribute ID")
	CrtCmd.Flags().StringVarP(&displayName, "display-name", "d",
		"", "Attribute Display Name")
	CrtCmd.Flags().StringVarP(&description, "description", "",
		"", "Attribute Description")
	CrtCmd.Flags().StringVarP(&scope, "scope", "s",
		"", "Attribute scope")
	CrtCmd.Flags().StringVarP(&dataType, "data-type", "",
		"", "Attribute data type")
	CrtCmd.Flags().IntVarP(&cardinality, "cardinality", "c",
		1, "Attribute cardinality")
	CrtCmd.Flags().StringVarP(&aValuesPath, "allowed-values", "",
		"", "Path to a file containing allowed values")

	_ = CrtCmd.MarkFlagRequired("id")
	_ = CrtCmd.MarkFlagRequired("display-name")
	_ = CrtCmd.MarkFlagRequired("scope")
	_ = CrtCmd.MarkFlagRequired("data-type")
}
