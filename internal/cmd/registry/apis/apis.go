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

package apis

import (
	"internal/cmd/registry/apis/artifacts"
	"internal/cmd/registry/apis/deployments"
	"internal/cmd/registry/apis/versions"

	"github.com/spf13/cobra"
)

// APICmd to manage apis
var APICmd = &cobra.Command{
	Use:   "apis",
	Short: "Manage APIs in Apigee Registry",
	Long:  "Manage APIs in Apigee Registry",
}

var (
	name                string
	labels, annotations map[string]string
)

func init() {
	APICmd.AddCommand(CreateAPICmd)
	APICmd.AddCommand(GetAPICmd)
	APICmd.AddCommand(DelAPICmd)
	APICmd.AddCommand(ListAPICmd)
	APICmd.AddCommand(deployments.APIDeploymentCmd)
	APICmd.AddCommand(versions.APIVersionCmd)
	APICmd.AddCommand(artifacts.APIArtifactCmd)
}
