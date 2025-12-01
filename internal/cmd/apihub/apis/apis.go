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

package apis

import (
	"internal/cmd/apihub/apis/versions"

	"github.com/spf13/cobra"
)

// ApisCmd to manage apis
var ApisCmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage Apigee API Hub APIs",
	Long:  "Manage Apigee API Hub APIs",
}

var org, region string

var examples = []string{
	`apigeecli apihub apis create -i $apiId -f ./test/api.json -r $region -o $project --default-token`,
}

func init() {
	ApisCmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ApisCmd.PersistentFlags().StringVarP(&region, "region", "r",
		"", "API Hub region name")

	ApisCmd.AddCommand(ListCmd)
	ApisCmd.AddCommand(CrtCmd)
	ApisCmd.AddCommand(GetCmd)
	ApisCmd.AddCommand(DelCmd)
	ApisCmd.AddCommand(UpdateCmd)
	ApisCmd.AddCommand(ExportCmd)
	ApisCmd.AddCommand(versions.ApiVersionsCmd)

	_ = ApisCmd.MarkFlagRequired("org")
	_ = ApisCmd.MarkFlagRequired("region")
}

func GetExample(i int) string {
	return examples[i]
}
