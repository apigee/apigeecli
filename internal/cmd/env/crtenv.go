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

package env

import (
	"encoding/json"
	"fmt"
	"internal/apiclient"
	"internal/client/env"
	"internal/client/operations"
	"net/url"
	"path/filepath"

	"github.com/spf13/cobra"
)

// CreateCmd to create env
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new environment",
	Long:  "Create a new environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		_, err = url.Parse(fwdProxyURI)
		if err != nil {
			return fmt.Errorf("invalid URI string for fwdProxyURI: %v", err)
		}
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		respBody, err := env.Create(envType, deploymentType, fwdProxyURI)
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

var deploymentType, fwdProxyURI, envType string

func init() {
	CreateCmd.Flags().StringVarP(&environment, "env", "e",
		"", "Apigee environment name")
	CreateCmd.Flags().StringVarP(&envType, "envtype", "",
		"", "Environment type - must be BASE, INTERMEDIATE or COMPREHENSIVE")
	CreateCmd.Flags().StringVarP(&deploymentType, "deptype", "d",
		"", "Deployment type - must be PROXY or ARCHIVE")
	CreateCmd.Flags().StringVarP(&fwdProxyURI, "fwdproxyuri", "f",
		"", "URL of the forward proxy to be applied to the runtime instances in this env")
	CreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the create to finish, with success or error")
	_ = CreateCmd.MarkFlagRequired("env")
}
