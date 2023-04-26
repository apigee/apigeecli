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

	"internal/apiclient"
)

var validMemberTypes = []string{"serviceAccount", "group", "user", "domain"}

// GetIAM
func GetIAM() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv()+":getIamPolicy")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// SetIAM
func SetIAM(memberName string, permission string, memberType string) (err error) {
	if !isValidMemberType(memberType) {
		return fmt.Errorf("Invalid memberType. Valid types are %v", validMemberTypes)
	}
	return apiclient.SetIAMPermission(memberName, permission, memberType)
}

// RemoveIAM
func RemoveIAM(memberName string, role string) (err error) {
	parts := strings.Split(memberName, ":")
	if !isValidMemberType(parts[0]) {
		return fmt.Errorf("Invalid memberType. Valid types are %v", validMemberTypes)
	}
	return apiclient.RemoveIAMPermission(memberName, role)
}

// TestIAM
func TestIAM(resource string, verb string) (respBody []byte, err error) {
	permission := "apigee." + resource + "." + verb
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv()+":testIamPermissions")
	payload := "{\"permissions\":[\"" + permission + "\"]}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

func isValidMemberType(memberType string) bool {
	for _, validMember := range validMemberTypes {
		if memberType == validMember {
			return true
		}
	}
	return false
}
