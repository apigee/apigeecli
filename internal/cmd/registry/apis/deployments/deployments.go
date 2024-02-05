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

package deployments

import (
	"github.com/spf13/cobra"
)

// APIDeploymentCmd to manage apis
var APIDeploymentCmd = &cobra.Command{
	Use:   "deployments",
	Short: "Manage API Deployments in Apigee Registry",
	Long:  "Manage API Deployments in Apigee Registry",
}

func init() {
	APIDeploymentCmd.AddCommand(ListDeployCmd)
	APIDeploymentCmd.AddCommand(CreateDeployCmd)
	APIDeploymentCmd.AddCommand(GetDeployCmd)
	APIDeploymentCmd.AddCommand(DelDeployCmd)
	APIDeploymentCmd.AddCommand(TagDeployCmd)
	APIDeploymentCmd.AddCommand(RollbackDeployCmd)
}
