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
	"internal/client/instances"
	"internal/client/operations"
	"path/filepath"

	"github.com/spf13/cobra"
)

// DeleteAttachCmd to create a new instance
var DeleteAttachCmd = &cobra.Command{
	Use:   "detach",
	Short: "Detach an environment from an instance",
	Long:  "Detach an environment from an instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		respBody, err := instances.DetachEnv(name)
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
	DeleteAttachCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	DeleteAttachCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	DeleteAttachCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deletion to finish, with success or error")

	_ = DeleteAttachCmd.MarkFlagRequired("name")
	_ = DeleteAttachCmd.MarkFlagRequired("env")
}
