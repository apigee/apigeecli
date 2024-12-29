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

package versions

import (
	"internal/cmd/apihub/apis/versions/specs"

	"github.com/spf13/cobra"
)

// ApiVersionsCmd to manage apis
var ApiVersionsCmd = &cobra.Command{
	Use:   "versions",
	Short: "Manage Apigee API Hub API Versions",
	Long:  "Manage Apigee API Hub API Versions",
}

var org, region string

var examples = []string{
	`apigeecli apihub apis versions create -i $apVeriId --api-id $apiId \
	-f ./test/api-ver.json -r $region -p $project --default-token`,
}

func init() {
	ApiVersionsCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ApiVersionsCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	ApiVersionsCmd.AddCommand(CrtCmd)
	ApiVersionsCmd.AddCommand(ListCmd)
	ApiVersionsCmd.AddCommand(GetCmd)
	ApiVersionsCmd.AddCommand(DelCmd)
	ApiVersionsCmd.AddCommand(UpdateCmd)
	ApiVersionsCmd.AddCommand(specs.SpecsCmd)

	_ = ApiVersionsCmd.MarkFlagRequired("org")
	_ = ApiVersionsCmd.MarkFlagRequired("region")
}

func GetExample(i int) string {
	return examples[i]
}
