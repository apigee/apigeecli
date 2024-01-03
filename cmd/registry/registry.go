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

package registry

import (
	"github.com/apigee/apigeecli/cmd/registry/apis"
	"github.com/apigee/apigeecli/cmd/registry/artifacts"
	"github.com/apigee/apigeecli/cmd/registry/instances"
	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// Cmd to manage apis
var Cmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage Apigee API Registry APIs",
	Long:  "Manage Apigee API Registry APIs",
}

func init() {
	Cmd.AddCommand(apis.APICmd)
	Cmd.AddCommand(artifacts.ArtifactCmd)
	Cmd.AddCommand(instances.InstanceCmd)

	Cmd.PersistentFlags().StringVarP(&utils.Org, "org", "o",
		"", "Apigee organization name")
	Cmd.PersistentFlags().StringVarP(&utils.Region, "region", "",
		"", "Region where Apigee Registry is provisioned")
	Cmd.PersistentFlags().StringVarP(&utils.ProjectID, "projectID", "p",
		"", "Apigee Registry Project ID; Use if Apigee Orgniazation not provisioned in this project")
}
