// Copyright 2020 Google LLC
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

package sharedflows

import (
	"github.com/spf13/cobra"
)

// CreateCmd to create sharedflow
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a Sharedflow in an Apigee Org",
	Long:  "Creates a Sharedflow in an Apigee Org",
}

func init() {
	CreateCmd.AddCommand(GhCreateCmd)
	CreateCmd.AddCommand(BundleCreateCmd)
}
