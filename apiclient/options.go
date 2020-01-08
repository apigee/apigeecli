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

import (
	"os"

	"github.com/srinandan/apigeecli/clilog"
)

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
	PrintOutput    bool   //prints output from http calls
}

var options = ApigeeClientOptions{SkipCache: false, SkipCheck: true, SkipLogInfo: true, PrintOutput: false}

//NewApigeeClient sets up options to invoke Apigee APIs
func NewApigeeClient(o ApigeeClientOptions) {
	options = o
	clilog.Init(o.SkipLogInfo)
}

//SetApigeeOrg sets the org variable
func SetApigeeOrg(org string) {
	options.Org = org
}

//GetApigeeOrg gets the org variable
func GetApigeeOrg() string {
	return options.Org
}

//GetApigeeOrgP gets the pointer to org variable
func GetApigeeOrgP() *string {
	return &options.Org
}

//SetApigeeEnv set the env variable
func SetApigeeEnv(env string) {
	options.Env = env
}

//GetApigeeEnv gets the env variable
func GetApigeeEnv() string {
	return options.Env
}

//GetApigeeEnvP gets the pointer to env variable
func GetApigeeEnvP() *string {
	return &options.Env
}

//SetApigeeToken sets the access token for use with Apigee API calls
func SetApigeeToken(token string) {
	options.Token = token
}

//GetApigeeToken get the access token value in client opts (does not generate it)
func GetApigeeToken() string {
	return options.Token
}

//GetApigeeToken get the access token pointer
func GetApigeeTokenP() *string {
	return &options.Token
}

//SetProjectID sets the project id
func SetProjectID(projectID string) {
	options.ProjectID = projectID
}

//GetProjectID gets the project id
func GetProjectID() string {
	return options.ProjectID
}

//GetProjectID gets the project id pointer
func GetProjectIDP() *string {
	return &options.ProjectID
}

//SetServiceAccount
func SetServiceAccount(serviceAccount string) {
	options.ServiceAccount = serviceAccount
}

//GetServiceAccount
func GetServiceAccount() string {
	if options.ServiceAccount == "" {
		envVar := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if envVar != "" {
			options.ServiceAccount = envVar
		}
	}
	return options.ServiceAccount
}

//GetServiceAccountP
func GetServiceAccountP() *string {
	return &options.ServiceAccount
}

//IsSkipCheck
func IsSkipCheck() bool {
	return options.SkipCheck
}

//SkipCheck
func SkipCheck() *bool {
	return &options.SkipCheck
}

//IsSkipCache
func IsSkipCache() bool {
	return options.SkipCache
}

//SkipCache
func SkipCache() *bool {
	return &options.SkipCache
}

//IsSkipLogInfo
func IsSkipLogInfo() bool {
	return options.SkipLogInfo
}

//SkipLogInfo
func SkipLogInfo() *bool {
	return &options.SkipLogInfo
}

//SetLogIngo
func LogInfo(l bool) {
	options.SkipLogInfo = l
}

//PrintOutput
func SetPrintOutput(output bool) {
	options.PrintOutput = output
}

//GetPrintOutput
func GetPrintOutput() bool {
	return options.PrintOutput
}
