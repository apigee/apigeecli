// Copyright 2024 Google LLC
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

package sites

import (
	"fmt"
	"internal/apiclient"
	"internal/client/sites"
	"os"

	"github.com/spf13/cobra"
)

// AudiencesCmd is the parent for portal audience subcommands
var AudiencesCmd = &cobra.Command{
	Use:   "audiences",
	Short: "Manage Apigee portal audiences",
	Long:  "Manage Apigee portal audiences and developer membership",
}

var (
	siteID     string
	audienceID string
	filePath   string
	devEmail   string
	audAction  string
)

// ListAudiencesCmd lists all audiences for a portal
var ListAudiencesCmd = &cobra.Command{
	Use:   "list",
	Short: "List all audiences for a portal site",
	Long:  "List all portal audiences for the given Apigee integrated developer portal site",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = sites.ListAudiences(siteID)
		return
	},
}

// GetAudienceCmd gets a portal audience by ID
var GetAudienceCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a portal audience",
	Long:  "Get a portal audience by ID from the given Apigee integrated developer portal site",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = sites.GetAudience(siteID, audienceID)
		return
	},
}

// CreateAudienceCmd creates a portal audience
var CreateAudienceCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a portal audience",
	Long:  "Create a portal audience from a JSON file for the given Apigee integrated developer portal site",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if filePath == "" {
			return fmt.Errorf("required flag \"file\" not set")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		_, err = sites.CreateAudience(siteID, audienceID, content)
		return
	},
}

// UpdateAudienceCmd updates a portal audience
var UpdateAudienceCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a portal audience",
	Long:  "Update (replace) a portal audience from a JSON file",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if filePath == "" {
			return fmt.Errorf("required flag \"file\" not set")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		_, err = sites.UpdateAudience(siteID, audienceID, content)
		return
	},
}

// DeleteAudienceCmd deletes a portal audience
var DeleteAudienceCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a portal audience",
	Long:  "Delete a portal audience from the given Apigee integrated developer portal site",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = sites.DeleteAudience(siteID, audienceID)
		return
	},
}

// ManageAudienceDevCmd adds or removes a developer from a portal audience
var ManageAudienceDevCmd = &cobra.Command{
	Use:   "managedev",
	Short: "Add or remove a developer from a portal audience",
	Long:  "Add or remove a developer from an Apigee portal audience by email address",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if audAction != "add" && audAction != "remove" {
			return fmt.Errorf("action must be \"add\" or \"remove\"")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true
		_, err = sites.ManageDeveloperAudience(siteID, audienceID, devEmail, audAction)
		return
	},
}

func init() {
	// list
	ListAudiencesCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	_ = ListAudiencesCmd.MarkFlagRequired("site")

	// get
	GetAudienceCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	GetAudienceCmd.Flags().StringVarP(&audienceID, "name", "n", "", "Audience ID")
	_ = GetAudienceCmd.MarkFlagRequired("site")
	_ = GetAudienceCmd.MarkFlagRequired("name")

	// create
	CreateAudienceCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	CreateAudienceCmd.Flags().StringVarP(&audienceID, "name", "n", "", "Audience ID (optional, auto-generated if not set)")
	CreateAudienceCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to a JSON file containing the audience definition")
	_ = CreateAudienceCmd.MarkFlagRequired("site")
	_ = CreateAudienceCmd.MarkFlagRequired("file")

	// update
	UpdateAudienceCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	UpdateAudienceCmd.Flags().StringVarP(&audienceID, "name", "n", "", "Audience ID")
	UpdateAudienceCmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to a JSON file containing the updated audience definition")
	_ = UpdateAudienceCmd.MarkFlagRequired("site")
	_ = UpdateAudienceCmd.MarkFlagRequired("name")
	_ = UpdateAudienceCmd.MarkFlagRequired("file")

	// delete
	DeleteAudienceCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	DeleteAudienceCmd.Flags().StringVarP(&audienceID, "name", "n", "", "Audience ID")
	_ = DeleteAudienceCmd.MarkFlagRequired("site")
	_ = DeleteAudienceCmd.MarkFlagRequired("name")

	// managedev
	ManageAudienceDevCmd.Flags().StringVarP(&siteID, "site", "s", "", "Portal site ID")
	ManageAudienceDevCmd.Flags().StringVarP(&audienceID, "name", "n", "", "Audience ID")
	ManageAudienceDevCmd.Flags().StringVarP(&devEmail, "developer", "d", "", "Developer email address")
	ManageAudienceDevCmd.Flags().StringVarP(&audAction, "action", "", "add", "Action: add or remove the developer from the audience")
	_ = ManageAudienceDevCmd.MarkFlagRequired("site")
	_ = ManageAudienceDevCmd.MarkFlagRequired("name")
	_ = ManageAudienceDevCmd.MarkFlagRequired("developer")

	AudiencesCmd.AddCommand(ListAudiencesCmd)
	AudiencesCmd.AddCommand(GetAudienceCmd)
	AudiencesCmd.AddCommand(CreateAudienceCmd)
	AudiencesCmd.AddCommand(UpdateAudienceCmd)
	AudiencesCmd.AddCommand(DeleteAudienceCmd)
	AudiencesCmd.AddCommand(ManageAudienceDevCmd)
}
