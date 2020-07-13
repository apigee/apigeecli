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

package envoy

import (
	"fmt"

	"github.com/spf13/cobra"
)

//GenJwtCmd to get org details
var GenJwtCmd = &cobra.Command{
	Use:   "gen-jwt",
	Short: "Generate JSON Web Tokens for Apigee Envoy Connector",
	Long:  "Generate JSON Web Token for Apigee Envoy Connector",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var token string
		if token, err = GenerateToken(folder, expiry); err != nil {
			return err
		}
		fmt.Printf("%s", token)
		return nil
	},
}

var folder string
var expiry int

func init() {

	GenJwtCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing remote-service.* files")
	GenJwtCmd.Flags().IntVarP(&expiry, "exp", "x",
		10, "expiry in minutes; default 10 mins")

	_ = GenJwtCmd.MarkFlagRequired("folder")
}
