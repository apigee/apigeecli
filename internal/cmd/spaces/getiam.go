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

package spaces

import (
	"internal/apiclient"
	"internal/client/spaces"

	"github.com/spf13/cobra"
)

// GetIamCmd to get env iam details
var GetIamCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the IAM policy on a Space",
	Long:  "Gets the IAM policy on a Space",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		// apiclient.SetApigeeSpace(space)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = spaces.GetIAM(space)
		return
	},
	Example: `Get IAM policy for a space: ` + GetExample(1),
}

func init() {
	GetIamCmd.Flags().StringVarP(&space, "space", "",
		"", "Name of the space")

	_ = GetIamCmd.MarkFlagRequired("space")
}
