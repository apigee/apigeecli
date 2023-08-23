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

	"internal/bundlegen/proxybundle"

	"internal/client/sharedflows"

	"github.com/spf13/cobra"
)

// CreateCmd to create shared flow
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a sharedflow in an Apigee Org",
	Long:  "Creates a sharedflow in an Apigee Org",
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
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if sfZip != "" {
			_, err = sharedflows.Create(name, sfZip)
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

			if err = proxybundle.GenerateArchiveBundle(sfFolder, sfBundlePath); err != nil {
				return err
			}
			if _, err = sharedflows.Create(name, sfBundlePath); err != nil {
				return err
			}
			return os.Remove(sfBundlePath)
		}
		return err
	},
}

var sfZip, sfFolder string

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Sharedflow name")
	CreateCmd.Flags().StringVarP(&sfZip, "sf-zip", "p",
		"", "Path to the Sharedflow bundle/zip file")
	CreateCmd.Flags().StringVarP(&sfFolder, "sf-folder", "f",
		"", "Path to the Sharedflow Bundle; ex: ./test/sharedflowbundle")

	_ = CreateCmd.MarkFlagRequired("name")
}
