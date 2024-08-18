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
	"os"
	"path"
	"strconv"
	"sync"
	"text/tabwriter"
	"time"

	"internal/apiclient"

	"internal/clilog"
)

const (
	proxy_dimension       = "apiproxy"
	proxy_deployment_type = "proxy_deployment_type"
	selection             = "sum(message_count)"
)

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

type apiCalls struct {
	count           int
	standardCount   int
	extensibleCount int
	sync.Mutex
}

var ApiCalls = &apiCalls{count: 0}

func TotalAPICallsInMonthAsync(dimension bool, environment string, month int, year int, envDetails bool, wg *sync.WaitGroup) {
	defer wg.Done()
	var totalApiCalls, totalExtensible, totalStandard int
	var err error

	if totalApiCalls, totalExtensible, totalStandard, err = TotalAPICallsInMonth(dimension, environment, month, year); err != nil {
		clilog.Error.Println(err)
		return
	}

	if dimension {
		ApiCalls.incrementExtensibleCount(totalExtensible)
		ApiCalls.incrementStandardCount(totalStandard)
		ApiCalls.incrementCount(totalApiCalls)
		if envDetails {
			w := tabwriter.NewWriter(os.Stdout, 32, 4, 0, ' ', 0)
			fmt.Fprintf(w, "%s\t%d/%d\t%d\t%d", environment, month, year, totalExtensible, totalStandard)
			fmt.Fprintln(w)
			w.Flush()

		}
	} else {
		ApiCalls.incrementCount(totalApiCalls)
		if envDetails {
			w := tabwriter.NewWriter(os.Stdout, 32, 4, 0, ' ', 0)
			fmt.Fprintf(w, "%s\t%d/%d\t%d", environment, month, year, totalApiCalls)
			fmt.Fprintln(w)
			w.Flush()

		}
	}
}

func TotalAPICallsInMonth(dimension bool, environment string, month int, year int) (totalApiCalls int, totalExtensible int, totalStandard int, err error) {
	environmentReport, err := getReport(proxy_deployment_type, environment, month, year)
	if err != nil {
		return -1, -1, -1, err
	}
	for _, e := range environmentReport.Environments {
		for _, d := range e.Dimensions {
			if dimension {
				if d.Name == "EXTENSIBLE" {
					for _, m := range d.Metrics {
						calls, _ := strconv.Atoi(m.Values[0])
						totalExtensible = totalExtensible + calls
					}
				} else if d.Name == "STANDARD" {
					for _, m := range d.Metrics {
						calls, _ := strconv.Atoi(m.Values[0])
						totalStandard = totalStandard + calls
					}
				} else if d.Name == "(not set)" {
					for _, m := range d.Metrics {
						calls, _ := strconv.Atoi(m.Values[0])
						totalApiCalls = totalApiCalls + calls
					}
				}
			} else {
				for _, m := range d.Metrics {
					calls, _ := strconv.Atoi(m.Values[0])
					totalApiCalls = totalApiCalls + calls
				}
			}
		}
	}
	return totalApiCalls, totalExtensible, totalStandard, nil
}

func getReport(dimension string, environment string, month int, year int) (r report, err error) {
	var respBody []byte

	// throttle API Calls
	getDefaultRate := apiclient.GetRate()
	apiclient.SetRate(apiclient.ApigeeAnalyticsAPI)
	defer apiclient.SetRate(getDefaultRate)

	u, _ := url.Parse(apiclient.GetApigeeBaseURL())

	timeRange := fmt.Sprintf("%d/01/%d 00:00~%d/%d/%d 23:59", month, year, month, daysIn(time.Month(month), year), year)

	q := u.Query()
	q.Set("select", selection)
	q.Set("timeRange", timeRange)

	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", environment, "stats", dimension)

	if respBody, err = apiclient.HttpClient(u.String()); err != nil {
		return r, err
	}

	if err = json.Unmarshal(respBody, &r); err != nil {
		return r, err
	}

	return r, nil
}

// GetCount
func (c *apiCalls) GetCount() int {
	return c.count
}

// GetExtensibleCount
func (c *apiCalls) GetExtensibleCount() int {
	return c.extensibleCount
}

// GetStandardCount
func (c *apiCalls) GetStandardCount() int {
	return c.standardCount
}

// ResetCount
func (c *apiCalls) ResetCount() {
	c.Lock()
	defer c.Unlock()
	c.count = 0
	c.extensibleCount = 0
	c.standardCount = 0
}

// daysIn returns the number of days in a month for a given year.
// source https://groups.google.com/g/golang-nuts/c/W-ezk71hioo
func daysIn(m time.Month, year int) int {
	// This is equivalent to time.daysIn(m, year).
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// incrementCount synchronizes counting
func (c *apiCalls) incrementCount(total int) {
	c.Lock()
	defer c.Unlock()
	c.count = c.count + total
}

// incrementStandardCount synchronizes counting
func (c *apiCalls) incrementStandardCount(total int) {
	c.Lock()
	defer c.Unlock()
	c.standardCount = c.standardCount + total
}

// incrementExtensibleCount synchronizes counting
func (c *apiCalls) incrementExtensibleCount(total int) {
	c.Lock()
	defer c.Unlock()
	c.extensibleCount = c.extensibleCount + total
}
