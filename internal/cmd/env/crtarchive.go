// Copyright 2021 Google LLC
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

package env

import (
	"encoding/json"
	"fmt"
	"internal/apiclient"
	"internal/bundlegen/proxybundle"
	"internal/client/env"
	"internal/client/operations"
	"internal/clilog"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type op struct {
	Name     string         `json:"name,omitempty"`
	Metadata metadata       `json:"metadata,omitempty"`
	Done     bool           `json:"done,omitempty"`
	Error    operationError `json:"error,omitempty"`
}

type operationError struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type metadata struct {
	Type               string `json:"@type,omitempty"`
	OperationType      string `json:"operationType,omitempty"`
	TargetResourceName string `json:"targetResourceName,omitempty"`
	State              string `json:"state,omitempty"`
}

// CreateArchiveCmd to create env archive
var CreateArchiveCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new revision of archive in the environment",
	Long:  "Create a new revision of archive in the environment",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if zipfile != "" && folder != "" {
			return fmt.Errorf("both zipfile and folder path cannot be passed")
		}
		apiclient.SetApigeeEnv(environment)
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if folder != "" {
			zipfile = name + ".zip"
			if err = proxybundle.GenerateArchive(folder, zipfile); err != nil {
				return err
			}
		}

		respBody, err := env.CreateArchive(name, zipfile)
		if wait {
			archiveResponse := op{}
			if err = json.Unmarshal(respBody, &archiveResponse); err != nil {
				return err
			}

			s := strings.Split(archiveResponse.Name, "/")
			operationID := s[len(s)-1]

			clilog.Info.Printf("Deployment operation id is %s\n", operationID)
			clilog.Info.Printf("Checking deployment status in %d seconds\n", interval)

			apiclient.DisableCmdPrintHttpResponse()

			stop := apiclient.Every(interval*time.Second, func(time.Time) bool {
				var respOpsBody []byte
				respMap := op{}
				if respOpsBody, err = operations.Get(operationID); err != nil {
					return true
				}

				if err = json.Unmarshal(respOpsBody, &respMap); err != nil {
					return true
				}

				if respMap.Metadata.State == "IN_PROGRESS" {
					clilog.Info.Printf("Archive deployment status is: %s. Waiting %d seconds.\n", respMap.Metadata.State, interval)
					return true
				} else if respMap.Metadata.State == "FINISHED" {
					if respMap.Error == (operationError{}) {
						clilog.Info.Println("Archive deployment completed with status: ", respMap.Metadata.State)
					} else {
						clilog.Info.Printf("Archive deployment failed with status: %s", respMap.Error.Message)
					}
				} else {
					clilog.Info.Printf("Unknown state %s", respMap.Metadata.State)
				}
				return false
			})

			<-stop

		}
		return err
	},
}

var (
	zipfile, folder string
	wait            bool
)

const interval = 10

func init() {
	CreateArchiveCmd.Flags().StringVarP(&name, "name", "n",
		"", "Archive name")
	CreateArchiveCmd.Flags().StringVarP(&zipfile, "zipfile", "z",
		"", "Archive Zip file")
	CreateArchiveCmd.Flags().StringVarP(&folder, "folder", "f",
		"", "Archive Folder")
	CreateArchiveCmd.Flags().BoolVarP(&wait, "wait", "w",
		false, "Wait for deployment")

	_ = CreateArchiveCmd.MarkFlagRequired("name")
}
