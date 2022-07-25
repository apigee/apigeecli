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

package env

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"time"

	"github.com/apigee/apigeecli/apiclient"
)

const proxy_dimension = "apiproxy"
const selection = "sum(message_count)"

type report struct {
	Environments []environmentreport `json:"environments,omitempty"`
}

type environmentreport struct {
	Name       string      `json:"name,omitempty"`
	Dimensions []dimension `json:"dimensions,omitempty"`
}

type dimension struct {
	Name    string   `json:"name,omitempty"`
	Metrics []metric `json:"metrics,omitempty"`
}

type metric struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

func TotalAPICallsInMonth(environment string, month int, year int) (total int, err error) {

	var apiCalls int
	var respBody []byte

	u, _ := url.Parse(apiclient.BaseURL)

	timeRange := fmt.Sprintf("%d/01/%d 00:00~%d/%d/%d 23:59", month, year, month, daysIn(time.Month(month), year), year)

	q := u.Query()
	q.Set("select", selection)
	q.Set("timeRange", timeRange)

	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", environment, "stats", proxy_dimension)

	if respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String()); err != nil {
		return -1, err
	}

	environmentReport := report{}

	if err = json.Unmarshal(respBody, &environmentReport); err != nil {
		return -1, err
	}

	for _, e := range environmentReport.Environments {
		for _, d := range e.Dimensions {
			for _, m := range d.Metrics {
				calls, _ := strconv.Atoi(m.Values[0])
				apiCalls = apiCalls + calls
			}
		}
	}

	return apiCalls, nil
}

// daysIn returns the number of days in a month for a given year.
//source https://groups.google.com/g/golang-nuts/c/W-ezk71hioo
func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
