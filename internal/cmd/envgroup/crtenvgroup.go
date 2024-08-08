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

package envgroup

import (
	"encoding/json"
	"internal/apiclient"
	"path/filepath"

	"internal/client/envgroups"
	"internal/client/operations"

	"github.com/spf13/cobra"
)

// CreateCmd to create a new product
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Environment Group",
	Long:  "Create an Environment Group",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		respBody, err := envgroups.Create(name, hostnames)
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

var wait bool

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Environment Group")
	CreateCmd.Flags().StringArrayVarP(&hostnames, "hosts", "d",
		[]string{}, "A list of hostnames")
	CreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the create to finish, with success or error")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("hosts")
}
