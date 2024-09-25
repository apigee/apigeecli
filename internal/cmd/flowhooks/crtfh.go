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

package flowhooks

import (
	"fmt"
	"internal/apiclient"
	"internal/client/flowhooks"
	"strconv"

	"github.com/spf13/cobra"
)

// CreateCmd to create flow hooks
var CreateCmd = &cobra.Command{
	Use:   "attach",
	Short: "Attach a flowhook",
	Long:  "Attach a flowhook",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetApigeeEnv(env)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if continueOnErrStr != "" {
			continueOnErrPtr = new(bool)
			*continueOnErrPtr, err = strconv.ParseBool(continueOnErrStr)
			if err != nil {
				return fmt.Errorf("continueOnErr should be a boolean value: %v", err)
			}
		}
		_, err = flowhooks.Attach(name, description, sharedflow, continueOnErrPtr)
		return
	},
}

var (
	continueOnErrPtr                          *bool
	description, sharedflow, continueOnErrStr string
)

func init() {
	CreateCmd.Flags().StringVarP(&name, "name", "n",
		"", "Flowhook point")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description for the flowhook")
	CreateCmd.Flags().StringVarP(&sharedflow, "sharedflow", "s",
		"", "Sharedflow name")
	CreateCmd.Flags().StringVarP(&continueOnErrStr, "continue", "c",
		"", "Continue on error")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("sharedflow")
}
