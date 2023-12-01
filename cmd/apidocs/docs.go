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

package apidocs

import (
	"github.com/spf13/cobra"
)

// DocCmd to manage apis
var DocCmd = &cobra.Command{
	Use:   "documentation",
	Short: "Manage the documentation for the specified catalog item",
	Long:  "Manage the documentation for the specified catalog item",
}

func init() {
	DocCmd.AddCommand(GetDocCmd)
	DocCmd.AddCommand(UpdateDocCmd)
}
