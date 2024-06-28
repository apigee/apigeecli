// Copyright 2024 Google LLC
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

package specs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"internal/apiclient"
	"internal/client/hub"
	"internal/clilog"

	"github.com/spf13/cobra"
)

// GetCmd to get a catalog items
var GetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a spec for an API Version",
	Long:  "Get a spec for an API Version",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if contents && fileContentPath == "" {
			return fmt.Errorf("The output flag must be set when contents is set to true")
		}
		if fileContentPath != "" && !contents {
			return fmt.Errorf("The output flag must not be set when contents is set to false")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		if contents {
			apiclient.DisableCmdPrintHttpResponse()

			respBody, err := hub.GetApiVersionsSpecContents(apiID, versionID, specID)
			if err != nil {
				return err
			}

			var specContent map[string]string
			var payload []byte

			err = json.Unmarshal(respBody, &specContent)
			if err != nil {
				return err
			}

			payload, err = base64.StdEncoding.DecodeString(specContent["contents"])
			if err != nil {
				return err
			}

			err = apiclient.WriteByteArrayToFile(fileContentPath, false, payload)
			if err == nil {
				clilog.Info.Printf("Contents of Spec %s written to %s\n", specID, fileContentPath)
			}
			return err
		}
		_, err = hub.GetApiVersionSpec(apiID, versionID, specID)
		return
	},
}

var fileContentPath string
var contents bool

func init() {
	GetCmd.Flags().StringVarP(&versionID, "version", "v",
		"", "API Version ID")
	GetCmd.Flags().StringVarP(&apiID, "api-id", "",
		"", "API ID")
	GetCmd.Flags().StringVarP(&specID, "spec-id", "s",
		"", "Spec ID")
	GetCmd.Flags().BoolVarP(&contents, "contents", "c",
		false, "Get contents")
	GetCmd.Flags().StringVarP(&fileContentPath, "output", "",
		"", "Path to a file to write the contents of the spec")

	_ = GetCmd.MarkFlagRequired("api-id")
	_ = GetCmd.MarkFlagRequired("version")
	_ = GetCmd.MarkFlagRequired("spec-id")
}
