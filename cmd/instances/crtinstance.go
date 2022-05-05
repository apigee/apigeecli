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
	"fmt"
	"regexp"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/instances"
	"github.com/apigee/apigeecli/client/orgs"
	"github.com/spf13/cobra"
)

//Cmd to create a new instance
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Instance",
	Long:  "Create an Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if billingType, err = orgs.GetOrgField("billingType"); err != nil {
			return err
		}

		if billingType != "EVALUATION" {
			re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/locations\/([a-zA-Z0-9_-]+)\/keyRings\/([a-zA-Z0-9_-]+)\/cryptoKeys\/([a-zA-Z0-9_-]+)`)
			ok := re.Match([]byte(diskEncryptionKeyName))
			if !ok {
				return fmt.Errorf("disk encryption key must be of the format projects/{project-id}/locations/{location}/keyRings/{test}/cryptoKeys/{cryptoKey}")
			}
		}

		if ipRange != "" {
			re := regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}($|\/(22))$`)
			ok := re.Match([]byte(ipRange))
			if !ok {
				return fmt.Errorf("ipRange must be a valid CIDR of range /22")
			}
		}

		_, err = instances.Create(name, location, diskEncryptionKeyName, ipRange)

		return
	},
}

var diskEncryptionKeyName, ipRange, billingType string

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	CreateCmd.Flags().StringVarP(&location, "location", "l",
		"", "Instance location")
	CreateCmd.Flags().StringVarP(&diskEncryptionKeyName, "diskenc", "d",
		"", "CloudKMS key name")
	CreateCmd.Flags().StringVarP(&ipRange, "iprange", "",
		"", "Peering CIDR Range; must be /22 range")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("location")
}
