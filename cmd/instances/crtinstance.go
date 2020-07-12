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
)

//Cmd to create a new instance
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an Instance",
	Long:  "Create an Instance",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeOrg(org)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/locations\/([a-zA-Z0-9_-]+)\/keyRings\/([a-zA-Z0-9_-]+)\/cryptoKeys\/([a-zA-Z0-9_-]+)`)
		ok := re.Match([]byte(diskEncryptionKeyName))
		if !ok {
			return fmt.Errorf("custom role must be of the format projects/{project-id}/locations/{location}/keyRings/{test}/cryptoKeys/{cryptoKey}")
		}
		_, err = instances.Create(name, location, diskEncryptionKeyName)
		return
	},
}

var diskEncryptionKeyName string

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Name of the Instance")
	CreateCmd.Flags().StringVarP(&location, "location", "l",
		"", "Instance location")
	CreateCmd.Flags().StringVarP(&diskEncryptionKeyName, "diskenc", "d",
		"", "Instance location")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("location")
	_ = CreateCmd.MarkFlagRequired("diskenc")
}
