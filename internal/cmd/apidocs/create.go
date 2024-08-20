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
	"strconv"

	"github.com/spf13/cobra"
)

// CreateCmd to create a catalog items
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new catalog item",
	Long:  "Create a new catalog item",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if published != "" {
			_, err = strconv.ParseBool(published)
			if err != nil {
				return fmt.Errorf("published must be a boolean value: %v", err)
			}
		} else {
			published = "false"
		}

		if anonAllowed != "" {
			_, err = strconv.ParseBool(anonAllowed)
			if err != nil {
				return fmt.Errorf("allow-anon must be a boolean value: %v", err)
			}
		} else {
			anonAllowed = "false"
		}

		if requireCallbackUrl != "" {
			_, err = strconv.ParseBool(requireCallbackUrl)
			if err != nil {
				return fmt.Errorf("require-callback-url must be a boolean value: %v", err)
			}
		} else {
			requireCallbackUrl = "false"
		}

		if siteid == "" {
			return fmt.Errorf("siteid is a mandatory parameter")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		_, err = apidocs.Create(siteid, title, description, published,
			anonAllowed, apiProductName, requireCallbackUrl, imageUrl, categoryIds)
		return
	},
}

var (
	title, description, published, anonAllowed, apiProductName string
	requireCallbackUrl, imageUrl                               string
	categoryIds                                                []string
)

func init() {
	CreateCmd.Flags().StringVarP(&title, "title", "l",
		"", "The user-facing name of the catalog item")
	CreateCmd.Flags().StringVarP(&description, "desc", "d",
		"", "Description of the catalog item")
	CreateCmd.Flags().StringVarP(&published, "published", "p",
		"false", "Denotes whether the catalog item is published to the portal or is in a draft state")
	CreateCmd.Flags().StringVarP(&anonAllowed, "allow-anon", "",
		"false", "Boolean flag that manages user access to the catalog item")
	CreateCmd.Flags().StringVarP(&apiProductName, "api-product", "",
		"", "The name field of the associated API product")
	CreateCmd.Flags().StringVarP(&requireCallbackUrl, "require-callback-url", "",
		"false", "Whether a callback URL is required when this catalog item's developer app is created")
	CreateCmd.Flags().StringVarP(&imageUrl, "image-url", "",
		"", "Location of the image used for the catalog item in the catalog")

	CreateCmd.Flags().StringArrayVarP(&categoryIds, "category-ids", "",
		nil, "The IDs of the API categories to which this catalog item belongs")

	_ = CreateCmd.MarkFlagRequired("name")
	_ = CreateCmd.MarkFlagRequired("title")
	_ = CreateCmd.MarkFlagRequired("apiProductName")
}
