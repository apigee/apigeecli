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
	"fmt"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

func CreateTraceOverrides(apiproxy string, exporter string, endpoint string, sampler string, sample_rate string) (respBody []byte, err error) {
	if sampler != "OFF" && sampler != "PROBABILITY" {
		return nil, fmt.Errorf("invalid sampler value. Must be OFF or PROBABILITY")
	}

	traceConfig := []string{}
	traceConfig = append(traceConfig, "\"apiProxy\":\""+apiproxy+"\"")
	traceConfig = append(traceConfig, getSamplingConfig(sampler, sample_rate))

	payload := "{" + strings.Join(traceConfig, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig", "overrides")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func GetTraceOverrides(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig", "overrides", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

func DeleteTraceOverrides(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig", "overrides", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

func ListTraceOverrides() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig", "overrides")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}
