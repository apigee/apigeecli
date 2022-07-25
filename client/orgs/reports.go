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
	"os"
	"text/tabwriter"
	"time"

	"github.com/apigee/apigeecli/apiclient"
	"github.com/apigee/apigeecli/client/env"
)

func TotalAPICallsInMonth(month int, year int, envDetails bool) (total int, err error) {

	var envListBytes []byte
	var envList []string
	var apiCalls int

	apiclient.SetPrintOutput(false)

	if envListBytes, err = env.List(); err != nil {
		return -1, err
	}

	if err = json.Unmarshal(envListBytes, &envList); err != nil {
		return -1, err
	}

	for _, environment := range envList {
		var calls int
		if calls, err = env.TotalAPICallsInMonth(environment, month, year); err != nil {
			return -1, err
		}
		if envDetails {
			w := tabwriter.NewWriter(os.Stdout, 26, 4, 0, ' ', 0)
			fmt.Fprintf(w, "%s\t%d/%d\t%d", environment, month, year, calls)
			fmt.Fprintln(w)
			w.Flush()
		}
		apiCalls = apiCalls + calls
	}

	apiclient.SetPrintOutput(true)

	return apiCalls, nil
}

func TotalAPICallsInYear(year int, envDetails bool) (total int, err error) {

	var monthlyTotal int

	t := time.Now()
	currentYear := t.Year()

	if year > currentYear {
		return -1, fmt.Errorf("Invalid year. Year cannot be greater than current year")
	}

	if currentYear == year {
		currentMonth := t.Month()
		for i := 1; i <= int(currentMonth); i++ {
			if monthlyTotal, err = TotalAPICallsInMonth(i, year, envDetails); err != nil {
				return -1, err
			}
			total = total + monthlyTotal
		}
	} else {
		for i := 1; i <= 12; i++ {
			if monthlyTotal, err = TotalAPICallsInMonth(i, year, envDetails); err != nil {
				return -1, err
			}
			total = total + monthlyTotal
		}
	}

	return total, nil
}
