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

package sharedflows

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"internal/apiclient"
	"internal/clilog"

	"internal/bundlegen/proxybundle"

	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// BundleCreateCmd to create shared flow
var BundleCreateCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Creates a sharedflow in an Apigee Org",
	Long:  "Creates a sharedflow in an Apigee Org; Optionally deploy the sharedflow to an env",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if sfZip != "" && sfFolder != "" {
			return fmt.Errorf("sharedflow bundle (zip) and folder to a sharedflow cannot be combined")
		}
		if sfZip == "" && sfFolder == "" {
			return fmt.Errorf("either sharedflow bundle (zip) or folder must be specified, not both")
		}
		if sfFolder != "" {
			if _, err := os.Stat(sfFolder); os.IsNotExist(err) {
				return err
			}
		}
		if env != "" {
			apiclient.SetApigeeEnv(env)
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var respBody []byte
		if sfZip != "" {
			respBody, err = sharedflows.Create(name, sfZip)
		} else if sfFolder != "" {
			if stat, err := os.Stat(folder); err == nil && !stat.IsDir() {
				return fmt.Errorf("supplied path is not a folder")
			}
			if filepath.Base(sfFolder) != "sharedflowbundle" {
				return fmt.Errorf("--sf-folder or -p must be a path to sharedflowbundle folder")
			}
			tmpDir, err := os.MkdirTemp("", "sf")
			if err != nil {
				return err
			}
			defer os.RemoveAll(tmpDir)

			sfBundlePath := path.Join(tmpDir, name+".zip")

			if err = proxybundle.GenerateArchiveBundle(sfFolder, sfBundlePath, true); err != nil {
				return err
			}
			if respBody, err = sharedflows.Create(name, sfBundlePath); err != nil {
				return err
			}
			if err = os.Remove(sfBundlePath); err != nil {
				return err
			}
		}
		if env != "" {
			clilog.Info.Printf("Deploying the Sharedflow %s to environment %s\n", name, env)
			if revision, err = GetRevision(respBody); err != nil {
				return err
			}
			if _, err = sharedflows.Deploy(name, revision, overrides, serviceAccountName); err != nil {
				return err
			}
			if wait {
				return Wait(name, revision)
			}
		}
		return err
	},
}

var sfZip, sfFolder string

func init() {
	BundleCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	BundleCreateCmd.Flags().StringVarP(&sfZip, "sf-zip", "p",
		"", "Path to the Sharedflow bundle/zip file")
	BundleCreateCmd.Flags().StringVarP(&sfFolder, "sf-folder", "f",
		"", "Path to the Sharedflow Bundle; ex: ./test/sharedflowbundle")
	BundleCreateCmd.Flags().StringVarP(&env, "env", "e",
		"", "Apigee environment name")
	BundleCreateCmd.Flags().BoolVarP(&overrides, "ovr", "",
		false, "Forces deployment of the new revision")
	BundleCreateCmd.Flags().BoolVarP(&wait, "wait", "",
		false, "Waits for the deployment to finish, with success or error")
	BundleCreateCmd.Flags().StringVarP(&serviceAccountName, "sa", "s",
		"", "The format must be {ACCOUNT_ID}@{PROJECT}.iam.gserviceaccount.com.")

	_ = BundleCreateCmd.MarkFlagRequired("name")
}
