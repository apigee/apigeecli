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

package versions

import (
	"internal/cmd/registry/apis/versions/artifacts"
	"internal/cmd/registry/apis/versions/specs"

	"github.com/spf13/cobra"
)

// APIVersionCmd to manage apis
var APIVersionCmd = &cobra.Command{
	Use:   "versions",
	Short: "Manage API Versions in Apigee Registry",
	Long:  "Manage API Versions in Apigee Registry",
}

func init() {
	APIVersionCmd.AddCommand(ListVersionCmd)
	APIVersionCmd.AddCommand(CreateVersionCmd)
	APIVersionCmd.AddCommand(GetVersionCmd)
	APIVersionCmd.AddCommand(DelVersionCmd)
	APIVersionCmd.AddCommand(artifacts.ArtifactCmd)
	APIVersionCmd.AddCommand(specs.SpecCmd)
}
