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

package envgroup

import (
	"internal/apiclient"
	"internal/client/envgroups"

	"github.com/spf13/cobra"
)

// ImpCmd to list products
var ImpCmd = &cobra.Command{
	Use:   "import",
	Short: "Import a file containing environment group definitions",
	Long:  "Import a file containing environment group definitions",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		err = envgroups.Import(envGrpFile)
		return
	},
}

var envGrpFile string

func init() {
	ImpCmd.Flags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")
	ImpCmd.Flags().StringVarP(&envGrpFile, "file", "f",
		"", "File containing environment group definitions")

	_ = ImpCmd.MarkFlagRequired("file")
}
