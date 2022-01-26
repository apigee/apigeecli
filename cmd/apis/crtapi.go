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
	"fmt"
	"os"
	"path"
	"regexp"

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

		if proxyZip != "" && proxyFolder != "" {
			return fmt.Errorf("Proxy bundle (zip) and folder to an API proxy cannot be combined.")
		}

		if useGitHub, err = gitHubValidations(); err != nil {
			return err
		}

		if proxyZip != "" && (oasFile != "" || oasURI != "" || useGitHub) {
			return fmt.Errorf("Importing a bundle (--proxy) cannot be combined with importing via an OAS file or GitHub import")
		}

		if proxyZip != "" && (gqlFile != "" || gqlURI != "" || useGitHub) {
			return fmt.Errorf("Importing a bundle (--proxy) cannot be combined with importing via an GraphQL file or GitHub import")
		}

		if oasFile != "" && oasURI != "" {
			return fmt.Errorf("Cannot combine importing an OAS through a file and URI")
		}

		if gqlFile != "" && gqlURI != "" {
			return fmt.Errorf("Cannot combine importing a GraphQL schema through a file and URI")
		}

		if useGitHub && (oasFile != "" || oasURI != "") {
			return fmt.Errorf("Cannot combine importing via OAS document and Github")
		}

		if useGitHub && (gqlFile != "" || gqlURI != "") {
			return fmt.Errorf("Cannot combine importing via GraphQL schema and Github")
		}

		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if useGitHub {
			if err = proxybundle.GitHubImportBundle(ghOwner, ghRepo, ghPath); err != nil {
				proxybundle.CleanUp()
				return err
			}
			_, err = apis.CreateProxy(name, bundleName)
			proxybundle.CleanUp()
			return err
		} else if proxyZip != "" {
			_, err = apis.CreateProxy(name, proxyZip)
		} else if proxyFolder != "" {
			curDir, _ := os.Getwd()
			if err = proxybundle.GenerateArchiveBundle(proxyFolder, path.Join(curDir, name+".zip")); err != nil {
				return err
			}
			_, err = apis.CreateProxy(name, name+".zip")
			_ = os.Remove(name + ".zip")
			return err
		} else if oasFile != "" || oasURI != "" {
			var content []byte
			var oasDocName string
			if oasFile != "" {
				oasDocName, content, err = bundle.LoadDocumentFromFile(oasFile, validateSpec, formatValidation)
			} else {
				oasDocName, content, err = bundle.LoadDocumentFromURI(oasURI, validateSpec, formatValidation)
			}
			if err != nil {
				return err
			}

			err = bundle.GenerateAPIProxyDefFromOAS(name, oasDocName, skipPolicy, addCORS)
			if err != nil {
				return err
			}

			err = proxybundle.GenerateAPIProxyBundle(name, string(content), oasDocName, "oas", skipPolicy, addCORS)
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

const bundleName = "apiproxy.zip"

var proxyZip, proxyFolder, oasFile, oasURI, gqlFile, gqlURI string
var ghOwner, ghRepo, ghPath string
var importProxy, validateSpec, skipPolicy, addCORS, useGitHub, formatValidation bool

func init() {

	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	CreateCmd.Flags().StringVarP(&proxyZip, "proxy-zip", "z",
		"", "API Proxy Bundle path")
	CreateCmd.Flags().StringVarP(&proxyFolder, "proxy", "p",
		"", "API Proxy folder path")
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
	CreateCmd.Flags().BoolVarP(&addCORS, "add-cors", "",
		false, "Add a CORS policy")
	CreateCmd.Flags().BoolVarP(&formatValidation, "formatValidation", "",
		true, "disables validation of schema type formats")
	CreateCmd.Flags().StringVarP(&ghOwner, "gh-owner", "",
		"", "The github organization or username. ex: In https://github.com/srinandan, srinandan is the user")
	CreateCmd.Flags().StringVarP(&ghRepo, "gh-repo", "",
		"", "The github repo name. ex: https://github.com/srinandan/sample-apps, sample-apps is the repo")
	CreateCmd.Flags().StringVarP(&ghPath, "gh-proxy-path", "",
		"", "The path in the repo to the apiproxy folder. ex: my-repo/apiproxy")

	_ = CreateCmd.MarkFlagRequired("name")
}

func gitHubValidations() (bool, error) {

	if ghOwner == "" && ghRepo == "" && ghPath == "" {
		return false, nil
	}

	if ghOwner == "" && (ghRepo != "" || ghPath != "") {
		return false, fmt.Errorf("GitHub Owner must be set along with GitHub Repo and GitHub path")
	}

	if ghRepo == "" && (ghOwner != "" || ghPath != "") {
		return false, fmt.Errorf("GitHub repo must be set along with GitHub owner and GitHub path")
	}

	if ghPath == "" && (ghRepo != "" || ghOwner != "") {
		return false, fmt.Errorf("GitHub path must be set along with GitHub Repo and GitHub owner")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		return false, fmt.Errorf("Github access token must be set with this feature")
	}

	//(\w+)?\/apiproxy$
	re := regexp.MustCompile(`(\w+)?\/apiproxy$`)
	if ok := re.Match([]byte(ghPath)); !ok {
		return false, fmt.Errorf("Github path must end with /apiproxy")
	}

	if ghOwner != "" {
		return true, nil
	}

	return false, nil
}
