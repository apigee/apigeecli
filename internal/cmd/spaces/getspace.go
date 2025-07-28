// Copyright 2025 Google LLC
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

package spaces

import (
	"fmt"
	"internal/apiclient"
	"internal/client/spaces"

	"github.com/spf13/cobra"
)

// GetCmd to get Space
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an Apigee Space",
	Long:  "Get an Apigee Space",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if name == "" {
			return fmt.Errorf("name cannot be empty")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if name != "" {
			_, err = spaces.Get(name)
			return err
		}
		return
	},
}

func init() {
	GetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the space")

	_ = GetCmd.MarkFlagRequired("name")
}
