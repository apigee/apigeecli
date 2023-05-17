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
	"sync"
	"time"

	"internal/apiclient"

	"internal/clilog"

	"internal/client/env"
)

func TotalAPICallsInMonth(month int, year int, envDetails bool, conn int) (total int, err error) {
	var pwg sync.WaitGroup
	var envListBytes []byte
	var envList []string

	// ensure the count is reset to zero before calculating the next set
	defer env.ApiCalls.ResetCount()

	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	if envListBytes, err = env.List(); err != nil {
		return -1, err
	}

	if err = json.Unmarshal(envListBytes, &envList); err != nil {
		return -1, err
	}

	numEntities := len(envList)
	clilog.Debug.Printf("Found %d environments\n", numEntities)
	clilog.Debug.Printf("Generate report with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	// ensure connections aren't greater than entities
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		clilog.Debug.Printf("Creating reports for a batch %d of environments\n", (i + 1))
		go batchReport(envList[start:end], month, year, envDetails, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		clilog.Debug.Printf("Creating reports for remaining %d environments\n", remaining)
		go batchReport(envList[start:numEntities], month, year, envDetails, &pwg)
		pwg.Wait()
	}

	return env.ApiCalls.GetCount(), nil
}

func batchReport(envList []string, month int, year int, envDetails bool, pwg *sync.WaitGroup) {
	defer pwg.Done()
	// batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(envList))

	for _, environment := range envList {
		go env.TotalAPICallsInMonthAsync(environment, month, year, envDetails, &bwg)
	}

	bwg.Wait()
}

func TotalAPICallsInYear(year int, envDetails bool, conn int) (total int, err error) {
	var monthlyTotal int

	t := time.Now()
	currentYear := t.Year()

	if year > currentYear {
		return -1, fmt.Errorf("Invalid year. Year cannot be greater than current year")
	}

	if currentYear == year {
		currentMonth := t.Month()
		for i := 1; i <= int(currentMonth); i++ { // run the loop only till the current month
			if monthlyTotal, err = TotalAPICallsInMonth(i, year, envDetails, conn); err != nil {
				return -1, err
			}
			total = total + monthlyTotal
		}
	} else {
		for i := 1; i <= 12; i++ { // run the loop for each month
			if monthlyTotal, err = TotalAPICallsInMonth(i, year, envDetails, conn); err != nil {
				return -1, err
			}
			total = total + monthlyTotal
		}
	}

	return total, nil
}
