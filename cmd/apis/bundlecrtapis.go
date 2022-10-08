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

	"github.com/apigee/apigeecli/apiclient"
	proxybundle "github.com/apigee/apigeecli/bundlegen/proxybundle"
	"github.com/apigee/apigeecli/client/apis"
	"github.com/spf13/cobra"
)

var BundleCreateCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Creates an API proxy from an Zip or folder",
	Long:  "Creates an API proxy from an Zip or folder",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if proxyZip != "" && proxyFolder != "" {
			return fmt.Errorf("proxy bundle (zip) and folder to an API proxy cannot be combined")
		}
		if proxyZip == "" && proxyFolder == "" {
			return fmt.Errorf("either proxy bundle (zip) or folder must be specified, not both")
		}
		if proxyZipFolder != "" {
			if _, err := os.Stat(proxyZipFolder); os.IsNotExist(err) {
				return err
			}
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var curDir string

		if proxyZip != "" {
			_, err = apis.CreateProxy(name, proxyZip)
		} else if proxyFolder != "" {
			if proxyZipFolder == "" {
				curDir, _ = os.Getwd()
			} else {
				curDir = proxyZipFolder
			}

			proxyBundlePath := path.Join(curDir, name+".zip")

			if err = proxybundle.GenerateArchiveBundle(proxyFolder, proxyBundlePath); err != nil {
				return err
			}
			if _, err = apis.CreateProxy(name, proxyBundlePath); err != nil {
				return err
			}

			err = os.Remove(proxyBundlePath)
		}
		return err
	},
}

var proxyZip, proxyFolder, proxyZipFolder string

func init() {
	BundleCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	BundleCreateCmd.Flags().StringVarP(&proxyZip, "proxy-zip", "z",
		"", "Path to the API Proxy bundle/zip file")
	BundleCreateCmd.Flags().StringVarP(&proxyFolder, "proxy", "p",
		"", "Path to the apiproxy folder")
	BundleCreateCmd.Flags().StringVarP(&proxyZipFolder, "proxy-zip-folder", "f",
		"", "Path to write the API Proxy Bundle")

	_ = BundleCreateCmd.MarkFlagRequired("name")
}
