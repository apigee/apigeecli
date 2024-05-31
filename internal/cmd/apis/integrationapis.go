// Copyright 2021 Google LLC
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

package apis

import (
	"os"

	"internal/apiclient"
	"internal/client/apis"

	"internal/bundlegen"
	"internal/bundlegen/proxybundle"

	"github.com/spf13/cobra"
)

var IntegrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Creates an API proxy template for Application Integration",
	Long:  "Creates an API proxy template for Application Integration with an API Trigger",
	Example: `Create an API Proxy template for Application Integration:
` + GetExample(4),
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		tmpDir, err := os.MkdirTemp("", "proxy")
		if err != nil {
			return err
		}

		defer os.RemoveAll(tmpDir)

		if err = bundlegen.GenerateIntegrationAPIProxy(name, desc, apitrigger); err != nil {
			return err
		}

		if err = proxybundle.GenerateIntegrationAPIProxyBundle(name, integration, apitrigger, true); err != nil {
			return err
		}

		if importProxy {
			_, err = apis.CreateProxy(name, tmpDir)
		}
		return err
	},
}

var integration, apitrigger string

func init() {
	IntegrationCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	IntegrationCmd.Flags().StringVarP(&desc, "desc", "d",
		"", "API Proxy description")
	IntegrationCmd.Flags().StringVarP(&integration, "integration", "i",
		"", "Integration name")
	IntegrationCmd.Flags().StringVarP(&apitrigger, "trigger", "",
		"", "API Trigger name; don't include 'api_trigger/'")
	IntegrationCmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation")

	_ = IntegrationCmd.MarkFlagRequired("name")
	_ = IntegrationCmd.MarkFlagRequired("integration")
	_ = IntegrationCmd.MarkFlagRequired("trigger")
}
