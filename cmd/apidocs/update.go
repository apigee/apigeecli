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
	"strconv"

	"internal/apiclient"
	"internal/client/apidocs"

	"github.com/spf13/cobra"
)

// UpdateCmd to create a catalog items
var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing catalog item",
	Long:  "Update an existing catalog item",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		_, err = strconv.ParseBool(published)
		if err != nil {
			return fmt.Errorf("published must be a boolean value: %v", err)
		}

		_, err = strconv.ParseBool(anonAllowed)
		if err != nil {
			return fmt.Errorf("allow-anon must be a boolean value: %v", err)
		}

		_, err = strconv.ParseBool(requireCallbackUrl)
		if err != nil {
			return fmt.Errorf("require-callback-url must be a boolean value: %v", err)
		}

		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = apidocs.Update(siteid, id, title, description, published,
			anonAllowed, apiProductName, requireCallbackUrl, imageUrl, categoryIds)
		return
	},
}

func init() {
	UpdateCmd.Flags().StringVarP(&name, "id", "i",
		"", "Catalog UUID")
	UpdateCmd.Flags().StringVarP(&title, "title", "l",
		"", "The user-facing name of the catalog item")
	UpdateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description of the catalog item")
	UpdateCmd.Flags().StringVarP(&published, "published", "p",
		"", "Denotes whether the catalog item is published to the portal or is in a draft state")
	UpdateCmd.Flags().StringVarP(&anonAllowed, "allow-anon", "",
		"", "Boolean flag that manages user access to the catalog item")
	UpdateCmd.Flags().StringVarP(&apiProductName, "api-product", "",
		"", "The name field of the associated API product")
	UpdateCmd.Flags().StringVarP(&requireCallbackUrl, "require-callback-url", "",
		"", "Whether a callback URL is required when this catalog item's developer app is created")
	UpdateCmd.Flags().StringVarP(&imageUrl, "image-url", "",
		"", "Location of the image used for the catalog item in the catalog")

	UpdateCmd.Flags().StringArrayVarP(&categoryIds, "category-ids", "",
		nil, "The IDs of the API categories to which this catalog item belongs")

	_ = UpdateCmd.MarkFlagRequired("id")
	_ = UpdateCmd.MarkFlagRequired("title")
	_ = UpdateCmd.MarkFlagRequired("apiProductName")
}
