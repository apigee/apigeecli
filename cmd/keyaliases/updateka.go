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

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// UpdateCmd to create key aliases
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a Key Alias from PEM a file",
	Long:  "Update a Key Alias from PEM a file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		if certFile != "" && !utils.FileExists(certFile) {
			return fmt.Errorf("certFile was not found")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = keyaliases.CreateOrUpdateKeyCert(keystoreName,
			name,
			true,
			ignoreExpiry,
			ignoreNewLine,
			certFile,
			"",
			password)

		return err
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&keystoreName, "key", "k",
		"", "Name of the key store")
	UpdateCmd.Flags().StringVarP(&name, "alias", "s",
		"", "Name of the key alias")

	UpdateCmd.Flags().BoolVarP(&ignoreExpiry, "exp", "x",
		false, "Ignore expiry validation")
	UpdateCmd.Flags().BoolVarP(&ignoreNewLine, "nl", "w",
		false, "Ignore new line in cert chain")
	UpdateCmd.Flags().StringVarP(&certFile, "cert-filepath", "",
		"", "Path to the X509 certificate in PEM format")

	_ = UpdateCmd.MarkFlagRequired("alias")
	_ = UpdateCmd.MarkFlagRequired("key")
}
