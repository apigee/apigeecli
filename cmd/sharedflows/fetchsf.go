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
	"internal/apiclient"

	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// FetCmd to download shared flow
var FetCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Returns a zip-formatted shared flow bundle ",
	Long:  "Returns a zip-formatted shared flow bundle of code and config files",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if revision == -1 {
			if revision, err = sharedflows.GetHighestSfRevision(name); err != nil {
				return err
			}
		}
		return sharedflows.Fetch(name, revision)
	},
}

func init() {
	FetCmd.Flags().StringVarP(&name, "name", "n",
		"", "Shared flow Bundle Name")
	FetCmd.Flags().IntVarP(&revision, "rev", "v",
		-1, "Shared flow revision. If not set, the highest revision is used")

	_ = FetCmd.MarkFlagRequired("name")
}
