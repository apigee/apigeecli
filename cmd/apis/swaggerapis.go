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
	"fmt"

	"github.com/apigee/apigeecli/apiclient"
	bundle "github.com/apigee/apigeecli/bundlegen"
	"github.com/spf13/cobra"
)

var SwaggerCreateCmd = &cobra.Command{
	Use:   "swagger",
	Short: "Creates an API proxy from a Swagger Spec",
	Long:  "Creates an API proxy from a Swagger Spec for Cloud Endpoints/API Gateway",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		//var content []byte
		var oasDocName string
		if swaggerFile != "" {
			if oasDocName, _, err = bundle.LoadSwaggerFromFile(swaggerFile, validateSpec); err != nil {
				return err
			}
		}
		fmt.Println("Passed file")
		//Generate the apiproxy struct
		err = bundle.GenerateAPIProxyFromSwagger(name,
			oasDocName,
			skipPolicy,
			addCORS)

		return err
	},
}

var swaggerFile string

func init() {
	SwaggerCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	SwaggerCreateCmd.Flags().StringVarP(&swaggerFile, "swaggerfile", "f",
		"", "Swagger Specification file")
	SwaggerCreateCmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation from spec")
	SwaggerCreateCmd.Flags().BoolVarP(&validateSpec, "validate", "",
		true, "Validate Spec before generating proxy")
	SwaggerCreateCmd.Flags().BoolVarP(&skipPolicy, "skip-policy", "",
		false, "Skip adding the OAS Validate policy")
	SwaggerCreateCmd.Flags().BoolVarP(&addCORS, "add-cors", "",
		false, "Add a CORS policy")

	_ = SwaggerCreateCmd.MarkFlagRequired("name")
	_ = SwaggerCreateCmd.MarkFlagRequired("swaggerFile")
}
