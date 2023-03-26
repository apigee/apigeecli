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
	"strings"

	"internal/apiclient"

	"internal/client/sync"

	"github.com/spf13/cobra"
)

//{"identities":["serviceAccount:srinandans-apigee@srinandans-apigee.iam.gserviceaccount.com"]}

// Cmd to set identities
var RemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove identity with access to control plane resources",
	Long:  "Remove identity with access to control plane resources",
	Args: func(cmd *cobra.Command, args []string) error {
		if !strings.Contains(identity, ".iam.gserviceaccount.com") {
			return fmt.Errorf("identity[0] must have .iam.gserviceaccount.com suffix"+
				" and should not be a Google managed service account: %s", identity)
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = sync.Remove(identity)
		return

	},
}

func init() {

	RemoveCmd.Flags().StringVarP(&identity, "ity", "i",
		"", "IAM Identity")

	_ = RemoveCmd.MarkFlagRequired("ity")
}
