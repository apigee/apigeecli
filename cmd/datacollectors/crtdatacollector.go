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

package datacollectors

import (
	"fmt"

	"internal/apiclient"

	"internal/client/datacollectors"

	"github.com/spf13/cobra"
)

// CreateCmd to create a new data collector
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Data Collector",
	Long:  "Create a Data Collector",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if !isValidType() {
			return fmt.Errorf("invalid collector type %s. Valid types are %s", collectorType, allowedTypes)
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = datacollectors.Create(name, description, collectorType)
		return
	},
}

var (
	description, collectorType string
	allowedTypes               = [6]string{"STRING", "INTEGER", "FLOAT", "LONG", "DOUBLE", "BOOLEAN"}
)

func init() {
	CreateCmd.Flags().StringVarP(&description, "description", "d",
		"", "Description of the data collector")
	CreateCmd.Flags().StringVarP(&collectorType, "type", "p",
		"", "Type of the data collector")
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Data Collector Name")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("type")
}

func isValidType() bool {
	for _, allowedType := range allowedTypes {
		if allowedType == collectorType {
			return true
		}
	}
	return false
}
