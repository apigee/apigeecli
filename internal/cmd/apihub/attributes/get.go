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
	"fmt"
	"internal/apiclient"
	"internal/client/hub"

	"github.com/spf13/cobra"
)

// GetCmd to get a catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get details for an Attribute",
	Long:  "Get details for an Attribute",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		var count int
		for _, v := range []bool{slo, env, dep, specType} {
			if v {
				count++
			}
		}
		if count > 1 {
			return fmt.Errorf("only one of --slo, --env, --deployments, --spec-types can be specified")
		}
		if attributeID != "" && count != 0 {
			return fmt.Errorf("attributeID cannot be mixed with other flags")
		}
		if attributeID == "" && count == 0 {
			return fmt.Errorf("attributeID or --slo, --env, --deployments, --spec-types must be specified")
		}

		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		if slo {
			attributeID = "system-slo"
		} else if env {
			attributeID = "system-environment"
		} else if dep {
			attributeID = "system-deployment-type"
		} else if specType {
			attributeID = "system-spec-type"
		}
		_, err = hub.GetAttribute(attributeID)
		return
	},
}

var slo, dep, env, specType bool

func init() {
	GetCmd.Flags().StringVarP(&attributeID, "id", "i",
		"", "Attribute ID")
	GetCmd.Flags().BoolVarP(&slo, "slo", "",
		false, "Get System SLO Attributes")
	GetCmd.Flags().BoolVarP(&env, "env", "",
		false, "Get System Environment Attributes")
	GetCmd.Flags().BoolVarP(&dep, "deployments", "",
		false, "Get System Deployment Attributes")
	GetCmd.Flags().BoolVarP(&specType, "spec-types", "",
		false, "Get System Deployment Attributes")
}
