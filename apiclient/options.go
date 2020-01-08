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

package apiclient

//BaseURL is the Apigee control plane endpoint
const BaseURL = "https://apigee.googleapis.com/v1/organizations/"

// ApigeeClientOptions is the base struct to hold all command arguments
type ApigeeClientOptions struct {
	Org            string //Apigee org
	Env            string //Apigee environment
	Token          string //Google OAuth access token
	ServiceAccount string //Google service account json
	ProjectID      string //GCP Project ID
	SkipLogInfo    bool   //LogInfo controls the log level
	SkipCheck      bool   //skip checking access token expiry
	SkipCache      bool   //skip writing access token to file
}

var options = ApigeeClientOptions{SkipCache: false, SkipCheck: true, SkipLogInfo: true}

func NewApigeeClient(o ApigeeClientOptions) {
	options = o
}

func SetApigeeOrg(org string) {
	options.Org = org
}

func GetApigeeOrg() string {
	return options.Org
}

func GetApigeeOrgP() *string {
	return &options.Org
}

func SetApigeeEnv(env string) {
	options.Env = env
}

func GetApigeeEnv() string {
	return options.Env
}

func GetApigeeEnvP() *string {
	return &options.Env
}

func SetApigeeToken(token string) {
	options.Token = token
}

func GetApigeeToken() string {
	return options.Token
}

func GetApigeeTokenP() *string {
	return &options.Token
}

func SetProjectID(projectID string) {
	options.ProjectID = projectID
}

func GetProjectID() string {
	return options.ProjectID
}

func GetProjectIDP() *string {
	return &options.ProjectID
}

func SetServiceAccount(serviceAccount string) {
	options.ServiceAccount = serviceAccount
}

func GetServiceAccount() string {
	return options.ServiceAccount
}

func GetServiceAccountP() *string {
	return &options.ServiceAccount
}

func IsSkipCheck() bool {
	return options.SkipCheck
}

func SkipCheck() *bool {
	return &options.SkipCheck
}

func IsSkipCache() bool {
	return options.SkipCache
}

func SkipCache() *bool {
	return &options.SkipCache
}

func IsSkipLogInfo() bool {
	return options.SkipLogInfo
}

func SkipLogInfo() *bool {
	return &options.SkipLogInfo
}
