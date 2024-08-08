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

package instances

import (
	"encoding/json"
	"internal/apiclient"
	"path/filepath"

	"internal/client/instances"
	"internal/client/operations"

	"github.com/spf13/cobra"
)

// ReserveNatCmd activates NAT for an Apigee instance
var ReserveNatCmd = &cobra.Command{
	Use:   "reserve",
	Short: "Reserve a NAT IP for an Apigeee instance",
	Long:  "Reserve a NAT IP for an Apigeee instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		respBody, err := instances.ReserveNatIP(name, natid)
		if err != nil {
			return
		}
		if wait {
			respMap := make(map[string]interface{})
			if err = json.Unmarshal(respBody, &respMap); err != nil {
				return err
			}
			err = operations.WaitForOperation(filepath.Base(respMap["name"].(string)))
		}
		return
	},
}

func init() {
	ReserveNatCmd.Flags().StringVarP(&natid, "natid", "i",
		"", "NAT identifier")
	ReserveNatCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the reserve to finish, with success or error")

	_ = ReserveNatCmd.MarkFlagRequired("natid")
}
