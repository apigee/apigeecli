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

package sites

import (
	"internal/apiclient"
	"internal/client/sites"

	"github.com/spf13/cobra"
)

// GetCmd gets an integrated developer portal by site ID
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get an Apigee integrated developer portal",
	Long:  "Get an Apigee integrated developer portal by its site ID",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = sites.Get(siteIDParam)
		return
	},
}

var siteIDParam string

func init() {
	GetCmd.Flags().StringVarP(&siteIDParam, "site", "s", "", "Portal site ID")
	_ = GetCmd.MarkFlagRequired("site")
}
