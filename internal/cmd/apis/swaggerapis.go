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
	"internal/apiclient"
	"internal/client/apis"

	bundle "internal/bundlegen"
	proxybundle "internal/bundlegen/proxybundle"

	"github.com/spf13/cobra"
)

var SwaggerCreateCmd = &cobra.Command{
	Use:   "swagger",
	Short: "Creates an API proxy from a Swagger Spec",
	Long:  "Creates an API proxy from a Swagger Spec for Cloud Endpoints/API Gateway",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if swaggerFile == "" && swaggerURI == "" {
			return fmt.Errorf("either swaggerfile or swaggeruri must be passed")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		// var content []byte
		var oasDocName string

		if swaggerURI != "" {
			if oasDocName, _, err = bundle.LoadSwaggerFromUri(swaggerURI); err != nil {
				return err
			}
		}

		if swaggerFile != "" {
			if oasDocName, _, err = bundle.LoadSwaggerFromFile(swaggerFile); err != nil {
				return err
			}
		}

		// Generate the apiproxy struct
		name, err = bundle.GenerateAPIProxyFromSwagger(name,
			desc,
			oasDocName,
			basePath,
			addCORS)
		if err != nil {
			return err
		}

		// Create the API proxy bundle
		err = proxybundle.GenerateAPIProxyBundleFromSwagger(name,
			skipPolicy,
			addCORS)

		if importProxy {
			_, err = apis.CreateProxy(name, name+zipExt, space)
		}

		return err
	},
}

var swaggerFile, swaggerURI string

func init() {
	SwaggerCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name. If not specified, will look for the x-google-api-name extension")
	SwaggerCreateCmd.Flags().StringVarP(&swaggerFile, "swaggerfile", "f",
		"", "Path to a Swagger Specification file with API Gateway or Cloud Endpoints extensions")
	SwaggerCreateCmd.Flags().StringVarP(&swaggerURI, "swaggeruri", "u",
		"", "URI to a Swagger Specification file with API Gateway or Cloud Endpoints extensions")
	SwaggerCreateCmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation from spec")
	SwaggerCreateCmd.Flags().StringVarP(&space, "space", "",
		"", "Apigee Space to associate to")
	SwaggerCreateCmd.Flags().BoolVarP(&addCORS, "add-cors", "",
		false, "Add a CORS policy")
	SwaggerCreateCmd.Flags().StringVarP(&desc, "desc", "d",
		"", "API Proxy description; Merges with the existing description in the spec.")
}
