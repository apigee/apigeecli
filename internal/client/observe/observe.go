// Copyright 2024 Google LLC
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

package observe

import (
	"net/url"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
)

type Action uint8

const (
	ENABLE Action = iota
	DISABLE
)

func CreateObservationJob(observationJobId string, sources []string) (respBody []byte, err error) {
	var payload string

	u, _ := url.Parse(apiclient.GetAPIObserveURL())
	u.Path = path.Join(u.Path, "observationJobs")
	q := u.Query()
	q.Set("observationJobId", observationJobId)
	u.RawQuery = q.Encode()

	if len(sources) != 0 {
		payload = "{ \"sourees\":[" + strings.Join(sources, ", ") + "]}"
	}
	respBody, err = apiclient.HttpClient(u.String(), payload, "POST")
	return respBody, err
}

func GetObservationJob(observationJobId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetAPIObserveURL())
	u.Path = path.Join(u.Path, "observationJobs", observationJobId)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteObservationJob(observationJobId string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetAPIObserveURL())
	u.Path = path.Join(u.Path, "observationJobs", observationJobId)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func GetApiObservation(name string, observationJob string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetAPIObserveURL())
	u.Path = path.Join(u.Path, "observationJobs", observationJob, "apiObservations", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func ListApiObservations(observationJob string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetAPIObserveURL())
	u.Path = path.Join(u.Path, "observationJobs", observationJob, "apiObservations")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}
