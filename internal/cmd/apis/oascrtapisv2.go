// Copyright 2023 Google LLC
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
	"os"
	"regexp"

	"internal/apiclient"
	"internal/clilog"

	bundle "internal/bundlegen"
	proxybundle "internal/bundlegen/proxybundle"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

var OasCreatev2Cmd = &cobra.Command{
	Use:     "openapi",
	Aliases: []string{"oas"},
	Short:   "Creates an API proxy from an OpenAPI Specification",
	Long:    "Creates an API proxy from an OpenAPI Specification",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if oasFile == "" && oasURI == "" {
			return fmt.Errorf("either oas-base-folderpath or oas-base-uri must be passed")
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
		if targetURL != "" && targetServerName != "" {
			return fmt.Errorf("targetURL and targetServerName cannot be set at the same time")
		}

		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var content []byte

		if oasFile != "" {
			err = checkFolder()
			if err != nil {
				return err
			}
		}

		content, err = bundle.LoadDocument(oasFile, oasURI, specName, validateSpec)
		if err != nil {
			return err
		}

		version := bundle.GetModelVersion()
		if version != "" {
			re := regexp.MustCompile(`3\.1\.[0-9]`)
			if re.MatchString(version) {
				clilog.Warning.Println("OpenAPI 3.1 detected. Skipping policy validation.")
				skipPolicy = true
			}
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
				TargetServerName:                targetServerName,
			},
		}

		// Generate the apiproxy struct
		err = bundle.GenerateAPIProxyDefFromOASv2(name,
			basePath,
			specName,
			skipPolicy,
			addCORS,
			targetOptions)

		if err != nil {
			return err
		}

		// Create the API proxy bundle
		err = proxybundle.GenerateAPIProxyBundleFromOAS(name,
			string(content),
			specName,
			skipPolicy,
			addCORS,
			targetOptions)

		if err != nil {
			return err
		}

		if importProxy {
			respBody, err := apis.CreateProxy(name, name+zipExt)
			if err != nil {
				return err
			}
			if env != "" {
				clilog.Info.Printf("Deploying the API Proxy %s to environment %s\n", name, env)
				if revision, err = GetRevision(respBody); err != nil {
					return err
				}
				if _, err = apis.DeployProxy(name, revision, overrides,
					sequencedRollout, safeDeploy, serviceAccountName); err != nil {
					return err
				}
				if wait {
					return apis.Wait(name, revision)
				}
			}
		}

		return err
	},
}

var (
	specName, oasFile, oasURI, targetURL, targetServerName                              string
	oasGoogleAcessTokenScopeLiteral, oasGoogleIDTokenAudLiteral, oasGoogleIDTokenAudRef string
	validateSpec, formatValidation                                                      bool
)

func init() {
	OasCreatev2Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	OasCreatev2Cmd.Flags().StringVarP(&basePath, "basepath", "p",
		"", "Base Path of the API Proxy; Overrides the basePath in spec")
	OasCreatev2Cmd.Flags().StringVarP(&oasFile, "oas-base-folderpath", "f",
		"", "Open API Spec Folder")
	OasCreatev2Cmd.Flags().StringVarP(&oasURI, "oas-base-uri", "u",
		"", "Open API Specification URI Base location")
	OasCreatev2Cmd.Flags().StringVarP(&specName, "oas-name", "",
		"", "Open API 3.0/3.1 Specification Name; Used with oas-base-filepath or oas-base-uri")
	OasCreatev2Cmd.Flags().StringVarP(&oasGoogleAcessTokenScopeLiteral, "google-accesstoken-scope-literal", "",
		"", "Generate Google Access token with target endpoint and set scope")
	OasCreatev2Cmd.Flags().StringVarP(&oasGoogleIDTokenAudLiteral, "google-idtoken-aud-literal", "",
		"", "Generate Google ID token with target endpoint and set audience")
	OasCreatev2Cmd.Flags().StringVarP(&oasGoogleIDTokenAudRef, "google-idtoken-aud-ref", "",
		"", "Generate Google ID token token with target endpoint and set audience reference")
	OasCreatev2Cmd.Flags().StringVarP(&targetURLRef, "target-url-ref", "",
		"", "Set a reference variable containing the target endpoint")
	OasCreatev2Cmd.Flags().StringVarP(&targetURL, "target-url", "",
		"", "Set a target URL for the target endpoint")
	OasCreatev2Cmd.Flags().StringVarP(&targetServerName, "target-server-name", "",
		"", "Set a target server name for the target endpoint")
	OasCreatev2Cmd.Flags().StringVarP(&integration, "integration", "i",
		"", "Integration name")
	OasCreatev2Cmd.Flags().StringVarP(&apitrigger, "trigger", "",
		"", "API Trigger name; don't include 'api_trigger/'")
	OasCreatev2Cmd.Flags().BoolVarP(&importProxy, "import", "",
		true, "Import API Proxy after generation from spec")
	OasCreatev2Cmd.Flags().BoolVarP(&validateSpec, "validate", "",
		false, "Validate Spec before generating proxy")
	OasCreatev2Cmd.Flags().BoolVarP(&skipPolicy, "skip-policy", "",
		false, "Skip adding the OAS Validate policy")
	OasCreatev2Cmd.Flags().BoolVarP(&addCORS, "add-cors", "",
		false, "Add a CORS policy")

	OasCreatev2Cmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	OasCreatev2Cmd.Flags().BoolVarP(&overrides, "ovr", "",
		false, "Forces deployment of the new revision")
	OasCreatev2Cmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deployment to finish, with success or error")
	OasCreatev2Cmd.Flags().BoolVarP(&sequencedRollout, "sequencedrollout", "",
		false, "If set to true, the routing rules will be rolled out in a safe order; default is false")
	OasCreatev2Cmd.Flags().BoolVarP(&safeDeploy, "safedeploy", "",
		true, deploymentMsg)
	OasCreatev2Cmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = OasCreatev2Cmd.MarkFlagRequired("name")
	_ = OasCreatev2Cmd.MarkFlagRequired("oas-name")
}

func checkFolder() error {
	f, err := os.Open(oasFile)
	if err != nil {
		return err
	}
	defer f.Close()
	fInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if !fInfo.IsDir() {
		return fmt.Errorf("oas-base-folderpath must be a folder")
	}
	return nil
}
