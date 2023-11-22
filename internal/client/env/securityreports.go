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

package env

import (
	"net/url"
	"path"
	"strconv"

	"internal/apiclient"
)

// GetSecurityReportView
func GetSecurityReportView(reportID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityReports", reportID, "resultView")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetSecurityReportResult
func GetSecurityReportResult(reportID string, name string) (err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityReports", reportID, "result")
	err = apiclient.DownloadResource(u.String(), name, ".zip", true)
	return err
}

// GetSecurityReport
func GetSecurityReport(reportID string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityReports", reportID)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListSecurityReports
func ListSecurityReports(pageSize int, pageToken string, dataset string, to string,
	from string, status string, submittedBy string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments",
		apiclient.GetApigeeEnv(), "securityReports")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	if dataset != "" {
		q.Set("dataset", dataset)
	}
	if to != "" {
		q.Set("to", to)
	}
	if from != "" {
		q.Set("from", from)
	}
	if status != "" {
		q.Set("status", status)
	}
	if submittedBy != "" {
		q.Set("submittedBy", submittedBy)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}
