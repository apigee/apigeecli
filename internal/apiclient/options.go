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
	"fmt"
	"os"
	"sync"

	"internal/clilog"
)

// BaseURL is the Apigee control plane endpoint
var (
	BaseURL      = "https://apigee.googleapis.com/v1/organizations/"
	StageBaseURL = "https://staging-apigee.sandbox.googleapis.com/v1/organizations/"
)

// ApigeeClientOptions is the base struct to hold all command arguments
type ApigeeClientOptions struct {
	Org            string // Apigee org
	Env            string // Apigee environment
	Token          string // Google OAuth access token
	ServiceAccount string // Google service account json
	ProjectID      string // GCP Project ID
	DebugLog       bool   // Enable debug logs
	TokenCheck     bool   // Check access token expiry
	SkipCache      bool   // skip writing access token to file
	PrintOutput    bool   // prints output from http calls
	NoOutput       bool   // Disable all statements to stdout
	ProxyUrl       string // use a proxy url
	MetadataToken  bool   // use metadata outh2 token
	APIRate        Rate   // throttle api calls to Apigee
}

var options *ApigeeClientOptions

type Rate uint8

const (
	None Rate = iota
	ApigeeAPI
	ApigeeAnalyticsAPI
)

var apiRate Rate

var cmdPrintHttpResponses = true

type clientPrintHttpResponse struct {
	enable bool
	sync.Mutex
}

var ClientPrintHttpResponse = &clientPrintHttpResponse{enable: true}

// NewApigeeClient sets up options to invoke Apigee APIs
func NewApigeeClient(o ApigeeClientOptions) {
	if options == nil {
		options = new(ApigeeClientOptions)
	}

	if o.Org != "" {
		options.Org = o.Org
	}
	if o.Token != "" {
		options.Token = o.Token
	}
	if o.ServiceAccount != "" {
		options.ServiceAccount = o.ServiceAccount
	}
	if o.ProjectID != "" {
		options.ProjectID = o.ProjectID
	}
	if o.Env != "" {
		options.Env = o.Env
	}

	options.TokenCheck = o.TokenCheck
	options.SkipCache = o.SkipCache
	options.DebugLog = o.DebugLog
	options.PrintOutput = o.PrintOutput
	options.NoOutput = o.NoOutput

	// initialize logs
	clilog.Init(options.DebugLog, options.PrintOutput, options.NoOutput)

	// read preference file
	_ = ReadPreferencesFile()
}

// UseStaging
func UseStaging() {
	BaseURL = StageBaseURL
}

// SetApigeeOrg sets the org variable
func SetApigeeOrg(org string) (err error) {
	if org == "" {
		if GetApigeeOrg() == "" {
			return fmt.Errorf("an org name was not set in preferences or supplied in the command")
		}
		return nil
	}
	options.Org = org
	return nil
}

// GetApigeeOrg gets the org variable
func GetApigeeOrg() string {
	return options.Org
}

// SetApigeeEnv set the env variable
func SetApigeeEnv(env string) {
	options.Env = env
}

// GetApigeeEnv gets the env variable
func GetApigeeEnv() string {
	return options.Env
}

// SetApigeeToken sets the access token for use with Apigee API calls
func SetApigeeToken(token string) {
	options.Token = token
}

// GetApigeeToken get the access token value in client opts (does not generate it)
func GetApigeeToken() string {
	return options.Token
}

// SetProjectID sets the project id
func SetProjectID(projectID string) {
	options.ProjectID = projectID
}

// GetProjectID gets the project id
func GetProjectID() string {
	return options.ProjectID
}

// SetServiceAccount
func SetServiceAccount(serviceAccount string) {
	options.ServiceAccount = serviceAccount
}

// GetServiceAccount
func GetServiceAccount() string {
	if options.ServiceAccount == "" {
		envVar := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if envVar != "" {
			options.ServiceAccount = envVar
		}
	}
	return options.ServiceAccount
}

// TokenCheckEnabled
func TokenCheckEnabled() bool {
	return options.TokenCheck
}

// IsSkipCache
func IsSkipCache() bool {
	return options.SkipCache
}

// DebugEnabled
func DebugEnabled() bool {
	return options.DebugLog
}

// PrintOutput
func SetPrintOutput(output bool) {
	options.PrintOutput = output
}

// GetPrintOutput
func GetPrintOutput() bool {
	return options.PrintOutput
}

// DisableCmdPrintHttpResponse
func DisableCmdPrintHttpResponse() {
	cmdPrintHttpResponses = false
}

// EnableCmdPrintHttpResponse
func EnableCmdPrintHttpResponse() {
	cmdPrintHttpResponses = true
}

// GetPrintHttpResponseSetting
func GetCmdPrintHttpResponseSetting() bool {
	return cmdPrintHttpResponses
}

// SetClientPrintHttpResponse
func (c *clientPrintHttpResponse) Set(b bool) {
	c.Lock()
	defer c.Unlock()
	c.enable = b
}

// GetPrintHttpResponseSetting
func (c *clientPrintHttpResponse) Get() bool {
	c.Lock()
	defer c.Unlock()
	return c.enable
}

// GetProxyURL
func GetProxyURL() string {
	return options.ProxyUrl
}

// SetProxyURL
func SetProxyURL(proxyurl string) {
	options.ProxyUrl = proxyurl
}

// DryRun
func DryRun() bool {
	if os.Getenv("APIGEECLI_DRYRUN") != "" {
		clilog.Warning.Println("Dry run mode enabled! unset APIGEECLI_DRYRUN to disable dry run")
		return true
	}
	return false
}

// SetNoOutput
func SetNoOutput(b bool) {
	options.NoOutput = b
}

// GetNoOutput
func GetNoOutput() bool {
	return options.NoOutput
}

// SetRate
func SetRate(r Rate) {
	apiRate = r
}

// GetRate
func GetRate() Rate {
	return apiRate
}
