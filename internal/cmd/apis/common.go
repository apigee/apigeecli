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
