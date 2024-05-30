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

package developers

import (
	"internal/apiclient"

	"internal/client/developers"

	"github.com/spf13/cobra"
)

// AdjustCmd to list developer subscriptions
var AdjustCmd = &cobra.Command{
	Use:   "adjust",
	Short: "Adjust the prepaid balance for the developer",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	Long: "Adjust the prepaid balance for the developer",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = developers.Adjust(email, adjust)
		return
	},
}

var adjust string

func init() {
	AdjustCmd.Flags().StringVarP(&email, "email", "n",
		"", "The developer's email")
	AdjustCmd.Flags().StringVarP(&adjust, "adjust", "",
		"", "Adjustment payload in json")

	_ = AdjustCmd.MarkFlagRequired("email")
	_ = AdjustCmd.MarkFlagRequired("adjust")
}
