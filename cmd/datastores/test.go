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

package datastores

import (
	"internal/apiclient"
	"internal/client/datastores"

	"github.com/spf13/cobra"
)

// TestCmd to update a datastore
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "Test an existing data store",
	Long:  "Test an existing data store",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if err = validate(targetType, bucketName, gcsPath,
			tablePrefix, datasetName); err != nil {
			return err
		}
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = datastores.Test(id, name, targetType, projectID, bucketName,
			gcsPath, datasetName, tablePrefix)
		return
	},
}

var id string

func init() {
	TestCmd.Flags().StringVarP(&id, "id", "i",
		"", "Data store id")
	TestCmd.Flags().StringVarP(&name, "name", "n",
		"", "Display name for the data store")
	TestCmd.Flags().StringVarP(&targetType, "target", "",
		"", "Destination storage type. Supported types gcs or bigquery")
	TestCmd.Flags().StringVarP(&projectID, "project-id", "p",
		"", "GCP project in which the datastore exists")
	TestCmd.Flags().StringVarP(&bucketName, "bucket-name", "b",
		"", "Name of the Cloud Storage bucket. Required for gcs targetType")
	TestCmd.Flags().StringVarP(&gcsPath, "gcs-path", "g",
		"", "Path of Cloud Storage bucket Required for gcs targetType")
	TestCmd.Flags().StringVarP(&datasetName, "dataset", "d",
		"", "BigQuery dataset name Required for bigquery targetType")
	TestCmd.Flags().StringVarP(&tablePrefix, "prefix", "f",
		"", "Prefix of BigQuery table Required for bigquery targetType")

	_ = TestCmd.MarkFlagRequired("id")
	_ = TestCmd.MarkFlagRequired("name")
	_ = TestCmd.MarkFlagRequired("target")
	_ = TestCmd.MarkFlagRequired("project-id")
}
