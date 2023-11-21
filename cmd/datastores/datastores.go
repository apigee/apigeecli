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
	"fmt"

	"github.com/spf13/cobra"
)

// Cmd to manage datastores
var Cmd = &cobra.Command{
	Use:   "datastores",
	Short: "Manage data stores to export analytics data",
	Long:  "Manage data stores to export analytics data",
}

var org, name, targetType string

func init() {
	Cmd.PersistentFlags().StringVarP(&org, "org", "o",
		"", "Apigee organization name")

	Cmd.AddCommand(ListCmd)
	Cmd.AddCommand(GetCmd)
	Cmd.AddCommand(DeleteCmd)
	Cmd.AddCommand(CreateCmd)
	Cmd.AddCommand(UpdateCmd)
	Cmd.AddCommand(TestCmd)

	_ = Cmd.MarkFlagRequired("org")
}

func validate(targetType string, bucketName string, gcsPath string,
	tablePrefix string, datasetName string,
) (err error) {
	if targetType != "gcs" && targetType != "bigquery" {
		return fmt.Errorf("invalid targetType; must be gcs or bigquery")
	}
	if targetType == "gcs" {
		if bucketName == "" || gcsPath == "" {
			return fmt.Errorf("bucketName and path are mandatory parameters for gcs targetType")
		}
	}
	if targetType == "bigquery" {
		if tablePrefix == "" || datasetName == "" {
			return fmt.Errorf("tablePrefix and datasetName are mandatory " +
				"parameters for bigquery targetType")
		}
	}
	return nil
}
