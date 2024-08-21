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

package keyaliases

import (
	"fmt"
	"internal/apiclient"
	"internal/client/keyaliases"
	"internal/clilog"
	"os"

	"github.com/spf13/cobra"
)

// GetctCmd to get key aliases
var GetctCmd = &cobra.Command{
	Use:   "getcert",
	Short: "Get a Key alias certificate",
	Long:  "Get a Key alias certificate",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if folder == "" {
			folder, err = os.Getwd()
			if err != nil {
				return err
			}
		}

		info, err := os.Stat(folder)
		if os.IsNotExist(err) {
			clilog.Error.Println("Directory does not exist")
			return err
		} else if err != nil {
			clilog.Error.Println("Error checking directory:", err)
			return err
		} else if !info.IsDir() {
			clilog.Error.Println("Directory does not exist")
			return fmt.Errorf("Directory does not exist")
		}
		err = keyaliases.GetCert(keystoreName, aliasName, folder)
		return
	},
}

var folder string

func init() {
	GetctCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	GetctCmd.Flags().StringVarP(&aliasName, "alias", "s",
		"", "Name of the key alias")
	GetctCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "Name of the folder to export the certificate")

	_ = GetctCmd.MarkFlagRequired("alias")
	_ = GetctCmd.MarkFlagRequired("key")
}
