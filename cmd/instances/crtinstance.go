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

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	"github.com/srinandan/apigeecli/client/instances"
	"github.com/srinandan/apigeecli/client/orgs"
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
				return fmt.Errorf("custom role must be of the format projects/{project-id}/locations/{location}/keyRings/{test}/cryptoKeys/{cryptoKey}")
			}

			if !isRangeValid(cidrRange) {
				return fmt.Errorf("Valid ranges are SLASH_16,SLASH_17,SLASH_18,SLASH_19 or SLASH_20")
			}
		}

		if billingType == "EVALUATION" {
			_, err = instances.Create(name, location, "", "SLASH_22")
		} else {
			_, err = instances.Create(name, location, diskEncryptionKeyName, cidrRange)
		}
		return
	},
}

var diskEncryptionKeyName, cidrRange, billingType string
var allowedRanges = []string{"SLASH_16", "SLASH_17", "SLASH_18", "SLASH_19", "SLASH_20"}

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	CreateCmd.Flags().StringVarP(&location, "location", "l",
		"", "Instance location")
	CreateCmd.Flags().StringVarP(&diskEncryptionKeyName, "diskenc", "d",
		"", "CloudKMS key name")
	CreateCmd.Flags().StringVarP(&cidrRange, "cidr", "r",
		"", "Peering CIDR Range; default is SLASH_16, other supported values SLASH_20, SLASH_22 for eval")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("location")
}

func isRangeValid(cidrRange string) bool {
	for _, allowedRange := range allowedRanges {
		if allowedRange == cidrRange {
			return true
		}
	}
	return false
}
