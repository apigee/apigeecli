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

package env

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

//GetTraceConfig
func GetTraceConfig() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//UpdateTraceConfig
func UpdateTraceConfig(exporter string, endpoint string, sampler string, sample_rate string) (respBody []byte, err error) {
	if exporter != "CLOUD_TRACE" && exporter != "JAEGER" {
		return nil, fmt.Errorf("invalid exporter format. Must be CLOUD_TRACE or JAEGER")
	}

	if sampler != "OFF" && sampler != "PROBABILITY" {
		return nil, fmt.Errorf("invalid sampler value. Must be OFF or PROBABILITY")
	}

	traceConfig := []string{}
	traceConfig = append(traceConfig, "\"exporter\":\""+exporter+"\"")
	if exporter == "JAEGER" {
		traceConfig = append(traceConfig, "\"endpoint\":\""+endpoint+"\"")
	}

	sampling := []string{}
	sampling = append(sampling, "\"sampler\":\""+sampler+"\"")
	sampling = append(sampling, "\"sampling_rate\":"+sample_rate)
	sampling_config := "\"sampling_config\":{" + strings.Join(sampling, ",") + "}"

	traceConfig = append(traceConfig, sampling_config)

	payload := "{" + strings.Join(traceConfig, ",") + "}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "traceConfig")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload, "PATCH")
	return respBody, err
}
