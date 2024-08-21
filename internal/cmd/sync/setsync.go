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

package sync

import (
	"fmt"
	"internal/apiclient"
	"internal/client/sync"
	"strings"

	"github.com/spf13/cobra"
)

// SetCmd to set identities
var SetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set identity with access to control plane resources",
	Long:  "Set identity with access to control plane resources",
	Args: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains(identity, ".iam.gserviceaccount.com") {
			return fmt.Errorf("identity[0] must have .iam.gserviceaccount.com suffix"+
				" and should not be a Google managed service account: %s", identity)
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = sync.Set(identity)
		return
	},
}

func init() {
	SetCmd.Flags().StringVarP(&identity, "ity", "i",
		"", "IAM Identity")

	_ = SetCmd.MarkFlagRequired("ity")
}
