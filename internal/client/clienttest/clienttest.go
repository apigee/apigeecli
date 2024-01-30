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

package clienttest

import (
	"fmt"
	"internal/apiclient"
	"os"
)

type TestRequirements uint8

const (
	ENV_REQD TestRequirements = iota
	ENV_NOT_REQD
	SITEID_REQD
	SITEID_NOT_REQD
	CLIPATH_REQD
	CLIPATH_NOT_REQD
)

func TestSetup(envReqd TestRequirements, siteIdReqd TestRequirements, cliPathReqd TestRequirements) (err error) {
	apiclient.NewApigeeClient(apiclient.ApigeeClientOptions{
		TokenCheck:  true,
		PrintOutput: true,
		NoOutput:    false,
		DebugLog:    true,
		SkipCache:   false,
	})

	org := os.Getenv("APIGEE_ORG")
	if err = apiclient.SetApigeeOrg(org); err != nil {
		return fmt.Errorf("APIGEE_ORG not set")
	}

	if envReqd == ENV_REQD {
		env := os.Getenv("APIGEE_ENV")
		if env == "" {
			return fmt.Errorf("APIGEE_ENV not set")
		}
		apiclient.SetApigeeEnv(env)
	}

	if siteIdReqd == SITEID_REQD {
		siteId := os.Getenv("APIGEE_SITEID")
		if siteId == "" {
			return fmt.Errorf("APIGEE_SITEID not set")
		}
	}

	region := os.Getenv("APIGEE_REGION")
	if region != "" {
		apiclient.SetRegion(region)
	}

	if apiclient.GetApigeeToken() == "" {
		token := os.Getenv("APIGEE_TOKEN")
		if token == "" {
			err = apiclient.GetDefaultAccessToken()
			if err != nil {
				return fmt.Errorf("APIGEE_TOKEN not set")
			}
		} else {
			apiclient.SetApigeeToken(token)
		}
	}

	if cliPathReqd == CLIPATH_REQD {
		cliPath := os.Getenv("APIGEECLI_PATH")
		if cliPath == "" {
			return fmt.Errorf("APIGEECLI_PATH not set")
		}
	}

	return nil
}
