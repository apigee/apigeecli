// Copyright 2023 Google LLC
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

package instances

import (
	"internal/apiclient"
	"internal/client/registry/instances"

	"internal/cmd/utils"

	"github.com/spf13/cobra"
)

// CreateInstanceCmd to create a new instance
var CreateInstanceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Apigee Registry Instance",
	Long:  "Create an Apigee Registry Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetProjectID(utils.ProjectID)
		apiclient.SetRegion(utils.Region)
		return apiclient.SetApigeeOrg(utils.Org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = instances.Create(cmekName)
		return err
	},
}

var cmekName string

func init() {
	CreateInstanceCmd.Flags().StringVarP(&cmekName, "cmek", "k",
		"", "Cloud KMS Key Resource ID")

	_ = CreateInstanceCmd.MarkFlagRequired("cmek")
}
