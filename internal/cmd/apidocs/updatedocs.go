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
	"fmt"
	"internal/apiclient"
	"internal/client/apidocs"
	"internal/cmd/utils"

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
		if graphQLPath != "" && endpointUri == "" {
			return fmt.Errorf("The flags graphQLPath and endpointUri must be set together")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		var openAPIDoc, graphQLDoc []byte

		if openAPIPath != "" {
			if openAPIDoc, err = utils.ReadFile(openAPIPath); err != nil {
				return err
			}
		}

		if graphQLPath != "" {
			graphQLDoc, err = utils.ReadFile(graphQLPath)
			if err != nil {
				return err
			}
		}

		_, err = apidocs.UpdateDocumentation(siteid, id, name, openAPIDoc, graphQLDoc, endpointUri)
		return
	},
}

var (
	openAPIPath, graphQLPath string
	endpointUri              string
)

func init() {
	UpdateDocCmd.Flags().StringVarP(&id, "id", "i",
		"", "Catalog UUID")
	UpdateDocCmd.Flags().StringVarP(&name, "name", "n",
		"", "Display Name")
	UpdateDocCmd.Flags().StringVarP(&openAPIPath, "openapi", "p",
		"", "Path to a file containing OpenAPI Specification documentation")
	UpdateDocCmd.Flags().StringVarP(&graphQLPath, "graphql", "q",
		"", "Path to a file containing GraphQL documentation")
	UpdateDocCmd.Flags().StringVarP(&endpointUri, "endpoint-uri", "e",
		"", "URI for the GraphQL proxy")

	_ = UpdateDocCmd.MarkFlagRequired("id")
	_ = UpdateDocCmd.MarkFlagRequired("name")
}
