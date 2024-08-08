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
	"fmt"
	"path/filepath"
	"regexp"

	"internal/apiclient"

	"internal/client/instances"
	"internal/client/operations"
	"internal/client/orgs"

	"github.com/spf13/cobra"
)

// CreateCmd to create a new instance
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Instance",
	Long:  "Create an Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		apiclient.ClientPrintHttpResponse.Set(false)
		if billingType, err = orgs.GetOrgField("billingType"); err != nil {
			return err
		}
		apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

		if billingType != "EVALUATION" {
			re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/locations` +
				`\/([a-zA-Z0-9_-]+)\/keyRings\/([a-zA-Z0-9_-]+)\/cryptoKeys\/([a-zA-Z0-9_-]+)`)
			ok := re.Match([]byte(diskEncryptionKeyName))
			if !ok {
				return fmt.Errorf("disk encryption key must be of the format " +
					"projects/{project-id}/locations/{location}/keyRings/{test}/cryptoKeys/{cryptoKey}")
			}
		}

		if ipRange != "" {
			re := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}($|\/(22))$`)
			ok := re.Match([]byte(ipRange))
			if !ok {
				return fmt.Errorf("ipRange must be a valid CIDR of range /22")
			}
		}

		respBody, err := instances.Create(name, location, diskEncryptionKeyName, ipRange, consumerAcceptList)
		if err != nil {
			return err
		}
		if wait {
			respMap := make(map[string]interface{})
			if err = json.Unmarshal(respBody, &respMap); err != nil {
				return err
			}
			err = operations.WaitForOperation(filepath.Base(respMap["name"].(string)))
		}

		return err
	},
}

var (
	diskEncryptionKeyName, ipRange, billingType string
	consumerAcceptList                          []string
)

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	CreateCmd.Flags().StringVarP(&location, "location", "l",
		"", "Instance location")
	CreateCmd.Flags().StringVarP(&diskEncryptionKeyName, "diskenc", "d",
		"", "CloudKMS key name")
	CreateCmd.Flags().StringVarP(&ipRange, "iprange", "",
		"", "Peering CIDR Range; must be /22 range")
	CreateCmd.Flags().StringArrayVarP(&consumerAcceptList, "consumer-accept-list", "c",
		[]string{}, "Customer accept list represents the list of "+
			"projects (id/number) that can connect to the service attachment")
	CreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the instance to finish, with success or error")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("location")
}
