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

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/apiclient"
	proxybundle "github.com/srinandan/apigeecli/bundlegen/proxybundle"
	"github.com/srinandan/apigeecli/client/apis"
)

var BundleCreateCmd = &cobra.Command{
	Use:   "bundle",
	Short: "Creates an API proxy from an Zip or folder",
	Long:  "Creates an API proxy from an Zip or folder",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if proxyZip != "" && proxyFolder != "" {
			return fmt.Errorf("Proxy bundle (zip) and folder to an API proxy cannot be combined.")
		}
		if proxyZip == "" && proxyFolder == "" {
			return fmt.Errorf("Either Proxy bundle (zip) or folder must be specified, not both")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if proxyZip != "" {
			_, err = apis.CreateProxy(name, proxyZip)
		} else if proxyFolder != "" {
			curDir, _ := os.Getwd()
			if err = proxybundle.GenerateArchiveBundle(proxyFolder, path.Join(curDir, name+".zip")); err != nil {
				return err
			}
			if _, err = apis.CreateProxy(name, name+".zip"); err != nil {
				return err
			}
			err = os.Remove(name + ".zip")
		}
		return err
	},
}

var proxyZip, proxyFolder string

func init() {
	BundleCreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	BundleCreateCmd.Flags().StringVarP(&proxyZip, "proxy-zip", "z",
		"", "API Proxy Bundle path")
	BundleCreateCmd.Flags().StringVarP(&proxyFolder, "proxy", "p",
		"", "API Proxy folder path")

	_ = BundleCreateCmd.MarkFlagRequired("name")
}
