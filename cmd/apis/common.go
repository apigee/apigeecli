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

package apis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"internal/apiclient"
	"internal/client/apis"
	"internal/clilog"
)

const deploymentMsg = "When set to true, generateDeployChangeReport will be executed and " +
	"deployment will proceed if there are no conflicts; default is true"

func GetRevision(respBody []byte) (revision int, err error) {
	var apiProxyRevResp map[string]interface{}

	err = json.Unmarshal(respBody, &apiProxyRevResp)
	if err != nil {
		return -1, err
	}
	apiProxyRev, err := strconv.Atoi(fmt.Sprintf("%v", apiProxyRevResp["revision"]))
	if err != nil {
		return -1, err
	}
	return apiProxyRev, nil
}

func Wait(name string, revision int) error {
	var err error

	clilog.Info.Printf("Checking deployment status in %d seconds\n", interval)

	apiclient.DisableCmdPrintHttpResponse()

	stop := apiclient.Every(interval*time.Second, func(time.Time) bool {
		var respBody []byte
		respMap := make(map[string]interface{})
		if respBody, err = apis.ListProxyRevisionDeployments(name, revision); err != nil {
			clilog.Error.Printf("Error fetching proxy revision status: %v", err)
			return false
		}

		if err = json.Unmarshal(respBody, &respMap); err != nil {
			return true
		}

		switch respMap["state"] {
		case "PROGRESSING":
			clilog.Info.Printf("Proxy deployment status is: %s. Waiting %d seconds.\n", respMap["state"], interval)
			return true
		case "READY":
			clilog.Info.Println("Proxy deployment completed with status: ", respMap["state"])
		default:
			clilog.Info.Println("Proxy deployment failed with status: ", respMap["state"])
		}

		return false
	})

	<-stop

	return err
}
