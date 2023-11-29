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
	"os"
	"path"
	"path/filepath"

	"internal/apiclient"
	"internal/clilog"

	proxybundle "internal/bundlegen/proxybundle"

	"internal/client/apis"

	"github.com/spf13/cobra"
)

var BundleCreateCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Creates an API proxy from an Zip or folder",
	Long:  "Creates an API proxy from an Zip or folder; Optionally deploy the API to an env",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if proxyZip != "" && proxyFolder != "" {
			return fmt.Errorf("proxy bundle (zip) and folder to an API proxy cannot be combined")
		}
		if proxyZip == "" && proxyFolder == "" {
			return fmt.Errorf("either proxy bundle (zip) or folder must be specified, not both")
		}
		if proxyFolder != "" {
			if _, err := os.Stat(proxyFolder); os.IsNotExist(err) {
				return err
			}
		}
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var respBody []byte
		if proxyZip != "" {
			respBody, err = apis.CreateProxy(name, proxyZip)
		} else if proxyFolder != "" {
			if stat, err := os.Stat(folder); err == nil && !stat.IsDir() {
				return fmt.Errorf("supplied path is not a folder")
			}
			if filepath.Base(proxyFolder) != "apiproxy" {
				return fmt.Errorf("--proxy-folder or -p must be a path to apiproxy folder")
			}
			tmpDir, err := os.MkdirTemp("", "proxy")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tmpDir)

			proxyBundlePath := path.Join(tmpDir, name+zipExt)

			if err = proxybundle.GenerateArchiveBundle(proxyFolder, proxyBundlePath, false); err != nil {
				return err
			}
			if respBody, err = apis.CreateProxy(name, proxyBundlePath); err != nil {
				return err
			}
			if err = os.Remove(proxyBundlePath); err != nil {
				return err
			}
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
				return Wait(name, revision)
			}
		}
		return err
	},
}

var proxyZip, proxyFolder string

func init() {
	BundleCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")

	BundleCreateCmd.Flags().StringVarP(&proxyZip, "proxy-zip", "p",
		"", "Path to the Proxy bundle/zip file")
	BundleCreateCmd.Flags().StringVarP(&proxyFolder, "proxy-folder", "f",
		"", "Path to the Proxy Bundle; ex: ./test/apiproxy")

	BundleCreateCmd.Flags().StringVarP(&env, "env", "e",
		"", "Name of the environment to deploy the proxy")
	BundleCreateCmd.Flags().BoolVarP(&overrides, "ovr", "r",
		false, "Forces deployment of the new revision")
	BundleCreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deployment to finish, with success or error")
	BundleCreateCmd.Flags().BoolVarP(&sequencedRollout, "sequencedrollout", "",
		false, "If set to true, the routing rules will be rolled out in a safe order; default is false")
	BundleCreateCmd.Flags().BoolVarP(&safeDeploy, "safedeploy", "",
		true, "When set to true, generateDeployChangeReport will be executed and "+
			"deployment will proceed if there are no conflicts; default is true")
	BundleCreateCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = BundleCreateCmd.MarkFlagRequired("name")
}
