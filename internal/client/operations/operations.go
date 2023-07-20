// Copyright 2020 Google LLC
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

package operations

import (
	"encoding/json"
	"net/url"
	"path"

	"internal/apiclient"
)

type ops struct {
	Operations []op `json:"operations,omitempty"`
}

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
	Type               string   `json:"@type,omitempty"`
	OperationType      string   `json:"operationType,omitempty"`
	TargetResourceName string   `json:"targetResourceName,omitempty"`
	State              string   `json:"state,omitempty"`
	Progress           progress `json:"progress,omitempty"`
}

type progress struct {
	Description string `json:"description,omitempty"`
	PercentDone int    `json:"percentDone,omitempty"`
}

type OperationCompleteState string

const (
	Success OperationCompleteState = "Success"
	Failed  OperationCompleteState = "Failed"
	Both    OperationCompleteState = "Both"
)

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "operations", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(state string, completeState OperationCompleteState) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "operations")
	if state != "" {
		apiclient.ClientPrintHttpResponse.Set(false)
		if respBody, err = apiclient.HttpClient(u.String()); err != nil {
			return nil, err
		}
		apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
		return filterOperation(respBody, state, completeState)
	}

	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func filterOperation(respBody []byte, state string, completeState OperationCompleteState) (operationsRespBody []byte, err error) {
	operationsResponse := ops{}
	filteredOperationsResponse := ops{}
	if err = json.Unmarshal(respBody, &operationsResponse); err != nil {
		return nil, err
	}

	for _, operationResponse := range operationsResponse.Operations {
		if operationResponse.Metadata.State == state {
			if completeState == Both {
				filteredOperationsResponse.Operations = append(filteredOperationsResponse.Operations, operationResponse)
			} else if completeState == Success {
				if operationResponse.Error == (operationError{}) {
					filteredOperationsResponse.Operations = append(filteredOperationsResponse.Operations, operationResponse)
				}
			} else if completeState == Failed {
				if operationResponse.Error != (operationError{}) {
					filteredOperationsResponse.Operations = append(filteredOperationsResponse.Operations, operationResponse)
				}
			}
		}
	}

	if operationsRespBody, err = json.Marshal(filteredOperationsResponse); err != nil {
		return nil, err
	}

	if apiclient.GetPrintOutput() {
		apiclient.PrettyPrint("json", operationsRespBody)
	}
	return operationsRespBody, nil
}
