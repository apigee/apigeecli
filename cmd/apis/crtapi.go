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

package apis

import (
	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	bundle "github.com/srinandan/apigeecli/bundlegen"
	proxybundle "github.com/srinandan/apigeecli/bundlegen/proxybundle"
	"github.com/srinandan/apigeecli/client/apis"
)

//CreateCmd to create api
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates an API proxy in an Apigee Org",
	Long:  "Creates an API proxy in an Apigee Org",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if proxy != "" {
			_, err = apis.CreateProxy(name, proxy)
		} else if oasFile != "" || oasURI != "" {
			var content []byte
			var oasDocName string
			if oasFile != "" {
				oasDocName, content, err = bundle.LoadDocumentFromFile(oasFile, validateSpec)
			} else {
				oasDocName, content, err = bundle.LoadDocumentFromURI(oasURI, validateSpec)
			}
			if err != nil {
				return err
			}

			err = bundle.GenerateAPIProxyDefFromOAS(name, oasDocName, skipPolicy)
			if err != nil {
				return err
			}

			err = proxybundle.GenerateAPIProxyBundle(name, string(content), oasDocName, skipPolicy)
			if err != nil {
				return err
			}

			if importProxy {
				_, err = apis.CreateProxy(name, name+".zip")
			}

		} else {
			_, err = apis.CreateProxy(name, "")
		}

		return
	},
}

var proxy, oasFile, oasURI string
var importProxy, validateSpec, skipPolicy bool

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	CreateCmd.Flags().StringVarP(&proxy, "proxy", "p",
		"", "API Proxy Bundle path")
	CreateCmd.Flags().StringVarP(&oasFile, "oasfile", "f",
		"", "Open API 3.0 Specification file")
	CreateCmd.Flags().StringVarP(&oasURI, "oasuri", "u",
		"", "Open API 3.0 Specification URI location")
	CreateCmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation from spec")
	CreateCmd.Flags().BoolVarP(&validateSpec, "validate", "",
		true, "Validate Spec before generating proxy")
	CreateCmd.Flags().BoolVarP(&skipPolicy, "skip-policy", "",
		false, "Skip adding the OAS Validate policy")

	_ = CreateCmd.MarkFlagRequired("name")
}
