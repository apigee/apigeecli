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

	bundle "internal/bundlegen"
	proxybundle "internal/bundlegen/proxybundle"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

var OasCreateCmd = &cobra.Command{
	Use:     "openapi",
	Aliases: []string{"oas"},
	Short:   "Creates an API proxy from an OpenAPI Specification",
	Long:    "Creates an API proxy from an OpenAPI Specification",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if oasFile == "" && oasURI == "" {
			return fmt.Errorf("either oasfile or oasuri must be passed")
		}
		if targetURL != "" && targetURLRef != "" {
			return fmt.Errorf("either target-url or target-url-ref must be passed, not both")
		}
		if integration != "" && apitrigger == "" {
			return fmt.Errorf("apitrigger must be passed if integration is set")
		}
		if integration == "" && apitrigger != "" {
			return fmt.Errorf("integration must be passed if apitrigger is set")
		}
		if (targetURL != "" || targetURLRef != "") && (integration != "" || apitrigger != "") {
			return fmt.Errorf("integration or apitrigger cannot be set if targetURL or targetURLRef is set")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var content []byte
		var oasDocName string

		integrationEndpoint := false

		if oasFile != "" {
			oasDocName, content, err = bundle.LoadDocumentFromFile(oasFile, validateSpec, formatValidation)
		} else {
			oasDocName, content, err = bundle.LoadDocumentFromURI(oasURI, validateSpec, formatValidation)
		}
		if err != nil {
			return err
		}

		targetOptions := bundle.TargetOptions{
			IntegrationBackend: bundle.IntegrationBackendOptions{
				IntegrationName: integration,
				TriggerName:     apitrigger,
			},
			HttpBackend: bundle.HttpBackendOptions{
				OasGoogleAcessTokenScopeLiteral: oasGoogleAcessTokenScopeLiteral,
				OasGoogleIDTokenAudLiteral:      oasGoogleIDTokenAudLiteral,
				OasGoogleIDTokenAudRef:          oasGoogleIDTokenAudRef,
				OasTargetURLRef:                 targetURLRef,
				TargetURL:                       targetURL,
			},
		}

		// check if integrationEndpoint is selected
		if integration != "" {
			integrationEndpoint = true
		}

		// Generate the apiproxy struct
		err = bundle.GenerateAPIProxyDefFromOAS(name,
			oasDocName,
			skipPolicy,
			addCORS,
			integrationEndpoint,
			oasGoogleAcessTokenScopeLiteral,
			oasGoogleIDTokenAudLiteral,
			oasGoogleIDTokenAudRef,
			targetURLRef,
			targetURL)

		if err != nil {
			return err
		}

		// Create the API proxy bundle
		err = proxybundle.GenerateAPIProxyBundleFromOAS(name,
			string(content),
			oasDocName,
			skipPolicy,
			addCORS,
			targetOptions)

		if err != nil {
			return err
		}

		if importProxy {
			_, err = apis.CreateProxy(name, name+".zip")
		}

		return err
	},
}

var (
	oasFile, oasURI, targetURL                                                          string
	oasGoogleAcessTokenScopeLiteral, oasGoogleIDTokenAudLiteral, oasGoogleIDTokenAudRef string
	validateSpec, formatValidation                                                      bool
)

func init() {
	OasCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	OasCreateCmd.Flags().StringVarP(&oasFile, "oasfile", "f",
		"", "Open API 3.0 Specification file")
	OasCreateCmd.Flags().StringVarP(&oasURI, "oasuri", "u",
		"", "Open API 3.0 Specification URI location")
	OasCreateCmd.Flags().StringVarP(&oasGoogleAcessTokenScopeLiteral, "google-accesstoken-scope-literal", "",
		"", "Generate Google Access token with target endpoint and set scope")
	OasCreateCmd.Flags().StringVarP(&oasGoogleIDTokenAudLiteral, "google-idtoken-aud-literal", "",
		"", "Generate Google ID token with target endpoint and set audience")
	OasCreateCmd.Flags().StringVarP(&oasGoogleIDTokenAudRef, "google-idtoken-aud-ref", "",
		"", "Generate Google ID token token with target endpoint and set audience reference")
	OasCreateCmd.Flags().StringVarP(&targetURLRef, "target-url-ref", "",
		"", "Set a reference variable containing the target endpoint")
	OasCreateCmd.Flags().StringVarP(&targetURL, "target-url", "",
		"", "Set a target URL for the target endpoint")
	OasCreateCmd.Flags().StringVarP(&integration, "integration", "i",
		"", "Integration name")
	OasCreateCmd.Flags().StringVarP(&apitrigger, "trigger", "",
		"", "API Trigger name; don't include 'api_trigger/'")
	OasCreateCmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation from spec")
	OasCreateCmd.Flags().BoolVarP(&validateSpec, "validate", "",
		true, "Validate Spec before generating proxy")
	OasCreateCmd.Flags().BoolVarP(&skipPolicy, "skip-policy", "",
		false, "Skip adding the OAS Validate policy")
	OasCreateCmd.Flags().BoolVarP(&addCORS, "add-cors", "",
		false, "Add a CORS policy")
	OasCreateCmd.Flags().BoolVarP(&formatValidation, "formatValidation", "",
		true, "disables validation of schema type formats")

	_ = OasCreateCmd.MarkFlagRequired("name")
}
