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
	"github.com/spf13/cobra"
)

//GenCmd to get org details
var GenCmd = &cobra.Command{
	Use:   "gen-jwk",
	Short: "Generate JSON Web Keys for Apigee Envoy Connector",
	Long:  "Generate JSON Web Keys for Apigee Envoy Connector",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = Generatekeys(kid)
		if err != nil {
			return err
		}
		return Generatekid(kid)
	},
}

var kid string

func init() {

	GenCmd.Flags().StringVarP(&kid, "kid", "k",
		"1", "Key Identifier")
}
