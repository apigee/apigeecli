// Copyright 2023 Google LLC
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

package apidocs

import (
	"encoding/json"
	"fmt"

	"internal/apiclient"
	"internal/client/apidocs"

	"github.com/apigee/apigeecli/cmd/utils"
	"github.com/spf13/cobra"
)

// UpdateDocCmd to get a catalog items
var UpdateDocCmd = &cobra.Command{
	Use:   "update",
	Short: "Updates the documentation for the specified catalog item",
	Long:  "Updates the documentation for the specified catalog item",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if openAPIPath == "" && graphQLPath == "" {
			return fmt.Errorf("The flags openapi and grapql cannot both be empty")
		}
		if openAPIPath != "" && graphQLPath != "" {
			return fmt.Errorf("The flags openapi and grapql cannot both be set")
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var openAPIDoc, graphQLDoc []byte
		var openAPI, graphQL string

		if openAPIPath != "" {
			if openAPIDoc, err = utils.ReadFile(openAPIPath); err != nil {
				return err
			}
			if openAPIDoc, err = json.Marshal(openAPIDoc); err != nil {
				return err
			}
			openAPI = string(openAPIDoc)
		}

		if graphQLPath != "" {
			graphQLDoc, err = utils.ReadFile(graphQLPath)
			if err != nil {
				return err
			}
			if graphQLDoc, err = json.Marshal(graphQLDoc); err != nil {
				return err
			}
			graphQL = string(graphQLDoc)
		}

		_, err = apidocs.UpdateDocumentation(siteid, name, openAPI, graphQL)
		return
	},
}

var openAPIPath, graphQLPath string

func init() {
	UpdateDocCmd.Flags().StringVarP(&name, "name", "n",
		"", "Catalog name")
	UpdateDocCmd.Flags().StringVarP(&openAPIPath, "openapi", "o",
		"", "Path to a file containing OpenAPI Specification documentation")
	UpdateDocCmd.Flags().StringVarP(&graphQLPath, "graphql", "q",
		"", "Path to a file containing GraphQL documentation")

	_ = UpdateDocCmd.MarkFlagRequired("name")
}
