// Copyright 2022 Google LLC
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

package orgs

import (
	"encoding/json"
	"fmt"
	"time"

	"internal/apiclient"

	"internal/clilog"

	"internal/client/env"
)

func TotalAPICallsInMonth(month int, year int, envDetails bool,
	proxyType bool, billingType string) (apiCalls int, extensibleApiCalls int, standardApiCalls int, err error) {

	var envListBytes []byte
	var envList []string

	// ensure the count is reset to zero before calculating the next set
	defer env.ApiCalls.ResetCount()

	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	if envListBytes, err = env.List(); err != nil {
		return -1, -1, -1, err
	}

	if err = json.Unmarshal(envListBytes, &envList); err != nil {
		return -1, -1, -1, err
	}

	numEntities := len(envList)
	clilog.Debug.Printf("Found %d environments\n", numEntities)

	batchReport(envList, month, year, envDetails, proxyType, billingType)

	return env.ApiCalls.GetCount(), env.ApiCalls.GetExtensibleCount(), env.ApiCalls.GetStandardCount(), nil
}

func batchReport(envList []string, month int, year int, envDetails bool, proxyType bool, billingType string) {
	for _, environment := range envList {
		env.TotalAPICallsInMonth(proxyType, environment, month, year, envDetails, billingType)
	}
}

func TotalAPICallsInYear(year int, envDetails bool,
	proxyType bool, billingType string) (apiCalls int, extensibleApiCalls int, standardApiCalls int, err error) {
	var monthlyApiTotal, monthlyExtensibleTotal, monthlyStandardTotal int

	t := time.Now()
	currentYear := t.Year()

	if year > currentYear {
		return -1, -1, -1, fmt.Errorf("Invalid year. Year cannot be greater than current year")
	}

	if currentYear == year {
		currentMonth := t.Month()
		for i := 1; i <= int(currentMonth); i++ { // run the loop only till the current month
			if monthlyApiTotal, monthlyExtensibleTotal, monthlyStandardTotal, err =
				TotalAPICallsInMonth(i, year, envDetails, proxyType, billingType); err != nil {
				return -1, -1, -1, err
			}
			apiCalls = apiCalls + monthlyApiTotal
			extensibleApiCalls = extensibleApiCalls + monthlyExtensibleTotal
			standardApiCalls = standardApiCalls + monthlyStandardTotal
		}
	} else {
		for i := 1; i <= 12; i++ { // run the loop for each month
			if monthlyApiTotal, monthlyExtensibleTotal, monthlyStandardTotal, err =
				TotalAPICallsInMonth(i, year, envDetails, proxyType, billingType); err != nil {
				return -1, -1, -1, err
			}
			apiCalls = apiCalls + monthlyApiTotal
			extensibleApiCalls = extensibleApiCalls + monthlyExtensibleTotal
			standardApiCalls = standardApiCalls + monthlyStandardTotal
		}
	}

	return apiCalls, extensibleApiCalls, standardApiCalls, nil
}
