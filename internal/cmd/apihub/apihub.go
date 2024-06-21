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

package apihub

import (
	"internal/cmd/apihub/apis"
	"internal/cmd/apihub/dependencies"
	"internal/cmd/apihub/deployments"
	"internal/cmd/apihub/externalapis"
	"internal/cmd/apihub/instances"

	"github.com/spf13/cobra"
)

// Cmd to manage apis
var Cmd = &cobra.Command{
	Use:   "apihub",
	Short: "Manage Apigee API Hub Resources",
	Long:  "Manage Apigee API Hub Resources",
}

func init() {
	Cmd.AddCommand(instances.InstanceCmd)
	Cmd.AddCommand(apis.ApisCmd)
	Cmd.AddCommand(deployments.DeploymentCmd)
	Cmd.AddCommand(dependencies.DependencyCmd)
	Cmd.AddCommand(externalapis.ExternalAPICmd)
}
