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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"

	"internal/clilog"
)

// CrmURL is the endpoint for cloud resource manager
const (
	CrmURL     = "https://cloudresourcemanager.googleapis.com/v1/projects/"
	crmBetaURL = "https://cloudresourcemanager.googleapis.com/v1beta1/projects/"
)

// binding for IAM Roles
type roleBinding struct {
	Role      string     `json:"role,omitempty"`
	Members   []string   `json:"members,omitempty"`
	Condition *condition `json:"condition,omitempty"`
}

// IamPolicy holds the response
type iamPolicy struct {
	Version  int           `json:"version,omitempty"`
	Etag     string        `json:"etag,omitempty"`
	Bindings []roleBinding `json:"bindings,omitempty"`
}

// SetIamPolicy holds the request to set IAM
type setIamPolicy struct {
	Policy iamPolicy `json:"policy,omitempty"`
}

// condition for Bindings
type condition struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

// CreateIAMServiceAccount create a new IAM SA with the necessary roles for Apigee
func CreateIAMServiceAccount(name string, iamRole string) (err error) {
	type KeyResponse struct {
		Name            string `json:"name,omitempty"`
		PrivateKeyType  string `json:"privateKeyType,omitempty"`
		PrivateKeyData  string `json:"privateKeyData,omitempty"`
		ValidBeforeTime string `json:"validBeforeTime,omitempty"`
		ValidAfterTime  string `json:"validAfterTime,omitempty"`
		KeyAlgorithm    string `json:"keyAlgorithm,omitempty"`
	}

	const iamURL = "https://iam.googleapis.com/v1/projects/"
	var role string

	serviceAccountName := name + "@" + GetProjectID() + ".iam.gserviceaccount.com"

	switch iamRole {
	case "sync":
		role = "roles/apigee.synchronizerManager"
	case "analytics":
		role = "roles/apigee.analyticsAgent"
	case "analyticsAgent":
		role = "roles/apigee.analyticsAgent"
	case "analyticsViewer":
		role = "roles/apigee.analyticsViewer"
	case "metric":
		role = "roles/monitoring.metricWriter"
	case "logger":
		role = "roles/logging.logWriter"
	case "mart":
		role = ""
	case "cassandra":
		role = "roles/storage.objectAdmin"
	case "connect":
		role = "roles/apigeeconnect.Agent"
	case "watcher":
		role = "roles/apigee.runtimeAgent"
	case "admin":
		role = "roles/apigee.admin"
	case "readonly-admin":
		role = "roles/apigee.readOnlyAdmin"
	case "api-admin":
		role = "roles/apigee.apiAdminV2"
	case "dev-admin":
		role = "roles/apigee.developerAdmin"
	case "env-admin":
		role = "roles/apigee.environmentAdmin"
	case "all":
		role = "not-necessary-to-add-this"
	default:
		return fmt.Errorf("invalid service account role")
	}

	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())

	// Step 1: create a new service account
	u, _ := url.Parse(iamURL)
	u.Path = path.Join(u.Path, GetProjectID(), "serviceAccounts")

	iamPayload := []string{}
	iamPayload = append(iamPayload, "\"accountId\":\""+name+"\"")
	iamPayload = append(iamPayload, "\"serviceAccount\": {\"displayName\": \""+name+"\"}")

	payload := "{" + strings.Join(iamPayload, ",") + "}"

	_, err = HttpClient(u.String(), payload)

	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	// Step 2: create a new service account key
	u, _ = url.Parse(iamURL)
	u.Path = path.Join(u.Path, GetProjectID(), "serviceAccounts",
		serviceAccountName, "keys")

	respKeyBody, err := HttpClient(u.String(), "")
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	// Step 3: read the response
	keyResponse := KeyResponse{}
	err = json.Unmarshal(respKeyBody, &keyResponse)
	if err != nil {
		return err
	}

	// Step 4: base64 decode the response to get the private key.json
	privateKey, err := base64.StdEncoding.DecodeString(keyResponse.PrivateKeyData)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	// Step 5: Write the data to a file
	file, err := os.Create(GetProjectID() + "-" + name + ".json")
	if err != nil {
		clilog.Error.Println("cannot open private key file: ", err)
		return err
	}

	defer file.Close()

	_, err = file.Write([]byte(privateKey))
	if err != nil {
		clilog.Error.Println("error writing to file: ", err)
		return err
	}

	// mart doesn't need any roles, return here.
	if iamRole == "mart" {
		return err
	}

	// Step 6: get the current IAM policies for the project
	u, _ = url.Parse(CrmURL)
	u.Path = path.Join(u.Path, GetProjectID()+":getIamPolicy")
	respBody, err := HttpClient(u.String(), "")

	iamPolicy := iamPolicy{}

	err = json.Unmarshal(respBody, &iamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	// Step 7: create a new policy binding for apigee
	if iamRole == "all" {
		bindings := createAllRoleBindings(serviceAccountName)
		iamPolicy.Bindings = append(iamPolicy.Bindings, bindings...)
	} else {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, "serviceAccount:"+serviceAccountName)

		iamPolicy.Bindings = append(iamPolicy.Bindings, binding)
	}

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = iamPolicy
	setIamPolicyBody, err := json.Marshal(setIamPolicy)

	// Step 8: set the iam policy
	u, _ = url.Parse(crmBetaURL)
	u.Path = path.Join(u.Path, GetProjectID()+":setIamPolicy")

	_, err = HttpClient(u.String(), string(setIamPolicyBody))

	return err
}

func createAllRoleBindings(name string) []roleBinding {
	roles := [...]string{
		"roles/apigee.synchronizerManager", "roles/apigee.analyticsAgent",
		"roles/monitoring.metricWriter", "roles/logging.logWriter", "roles/storage.objectAdmin",
		"roles/apigeeconnect.Agent", "roles/apigee.runtimeAgent",
	}

	bindings := []roleBinding{}

	for _, role := range roles {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, "serviceAccount:"+name)
		bindings = append(bindings, binding)
	}

	return bindings
}

// SetIAMPermission set permissions for a member on an Apigee Env
func SetIAMPermission(memberName string, iamRole string, memberType string) (err error) {
	var role string

	switch iamRole {
	case "sync":
		role = "roles/apigee.synchronizerManager"
	case "analytics":
		role = "roles/apigee.analyticsAgent"
	case "analyticsViewer":
		role = "roles/apigee.analyticsViewer"
	case "analyticsAgent":
		role = "roles/apigee.analyticsAgent"
	case "deploy":
		role = "roles/apigee.deployer"
	case "admin":
		role = "roles/apigee.environmentAdmin"
	default: // assume this is a custom role definition
		re := regexp.MustCompile(`projects\/([a-zA-Z0-9_-]+)\/roles\/([a-zA-Z0-9_-]+)`)
		result := re.FindString(iamRole)
		if result == "" {
			return fmt.Errorf("custom role must be of the format projects/{project-id}/roles/{role-name}")
		}
		role = iamRole
	}

	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())

	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":getIamPolicy")
	getIamPolicyBody, err := HttpClient(u.String())
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	getIamPolicy := iamPolicy{}

	err = json.Unmarshal(getIamPolicyBody, &getIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	foundRole := false
	for i, binding := range getIamPolicy.Bindings {
		if binding.Role == role {
			// found members with the role already, add the new SA to the role
			getIamPolicy.Bindings[i].Members = append(binding.Members, memberType+":"+memberName)
			foundRole = true
		}
	}

	// no members with the role, add a new one
	if !foundRole {
		binding := roleBinding{}
		binding.Role = role
		binding.Members = append(binding.Members, memberType+":"+memberName)
		getIamPolicy.Bindings = append(getIamPolicy.Bindings, binding)
	}

	u, _ = url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":setIamPolicy")

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = getIamPolicy

	setIamPolicyBody, err := json.Marshal(setIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	_, err = HttpClient(u.String(), string(setIamPolicyBody))

	return err
}

// RemoveIAMPermission removes/unbinds IAM permission from all roles for an Apigee Env
func RemoveIAMPermission(memberName string, iamRole string) (err error) {
	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())

	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":getIamPolicy")
	getIamPolicyBody, err := HttpClient(u.String())
	clilog.Debug.Println(string(getIamPolicyBody))
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	getIamPolicy := iamPolicy{}

	err = json.Unmarshal(getIamPolicyBody, &getIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	foundRole := false
	foundMember := false
	removeIamPolicy := setIamPolicy{}

	numBindings := len(getIamPolicy.Bindings)

	if numBindings < 1 {
		return fmt.Errorf("role %s not found for environment %s", iamRole, GetApigeeEnv())
	} else if numBindings == 1 { // there is only 1 binding
		clilog.Debug.Printf("comparing %s and %s\n", getIamPolicy.Bindings[0].Role, iamRole)
		if getIamPolicy.Bindings[0].Role == iamRole {
			if len(getIamPolicy.Bindings[0].Members) > 1 { // more than one member in the role
				removeIamPolicy.Policy.Etag = getIamPolicy.Etag
				// create a new role binding
				removeIamPolicy.Policy.Bindings = append(removeIamPolicy.Policy.Bindings, roleBinding{})
				// copy the role
				removeIamPolicy.Policy.Bindings[0].Role = getIamPolicy.Bindings[0].Role
				// copy other members
				for _, member := range getIamPolicy.Bindings[0].Members {
					clilog.Debug.Printf("comparing %s and %s\n", memberName, member)
					if member == memberName {
						clilog.Debug.Println("found member")
						foundMember = true
						// don't include this member
					} else {
						removeIamPolicy.Policy.Bindings[0].Members = append(removeIamPolicy.Policy.Bindings[0].Members, member)
					}
				}
				if !foundMember {
					return fmt.Errorf("member %s not set for role %s in environment %s", memberName, iamRole, GetApigeeEnv())
				}
			} else { // there is one member, one role
				if getIamPolicy.Bindings[0].Members[0] == memberName {
					clilog.Debug.Printf("comparing %s and %s\n", getIamPolicy.Bindings[0].Members[0], memberName)
					removeIamPolicy.Policy.Etag = getIamPolicy.Etag
				} else {
					return fmt.Errorf("member %s not set for role %s in environment %s", memberName, iamRole, GetApigeeEnv())
				}
			}
		} else {
			return fmt.Errorf("role %s not found for environment %s", iamRole, GetApigeeEnv())
		}
	} else { // there are many bindings, loop through them
		removeIamPolicy.Policy.Etag = getIamPolicy.Etag
		for _, binding := range getIamPolicy.Bindings {
			members := []string{}
			clilog.Debug.Printf("comparing %s and %s\n", binding.Role, iamRole)
			if binding.Role == iamRole {
				if len(binding.Members) > 1 { // there is more than one member in the role
					for _, member := range binding.Members {
						clilog.Debug.Printf("comparing %s and %s\n", member, memberName)
						if member == memberName { // remove the member
							foundMember = true
						} else {
							members = append(members, member)
						}
					}
					if !foundMember {
						return fmt.Errorf("member %s not set for role %s in environment %s", memberName, iamRole, GetApigeeEnv())
					}
				} else { // there is only one member in the role
					if binding.Members[0] == memberName {
						foundMember = true
					} else {
						return fmt.Errorf("member %s not set for role %s in environment %s", memberName, iamRole, GetApigeeEnv())
					}
				}
				copyRoleBinding := roleBinding{}
				copyRoleBinding.Role = binding.Role
				copyRoleBinding.Members = members
				removeIamPolicy.Policy.Bindings = append(removeIamPolicy.Policy.Bindings, copyRoleBinding)
				foundRole = true
			} else { // copy the binding as-is
				removeIamPolicy.Policy.Bindings = append(removeIamPolicy.Policy.Bindings, binding)
			}
		}

		if !foundRole {
			return fmt.Errorf("role %s not found for environment %s", iamRole, GetApigeeEnv())
		}
	}

	u, _ = url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), "environments", GetApigeeEnv()+":setIamPolicy")

	removeIamPolicyBody, err := json.Marshal(removeIamPolicy)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	_, err = HttpClient(u.String(), string(removeIamPolicyBody))
	return err
}

// AddWid add workload identity role to a service account
func AddWid(projectID string, namespace string, kServiceAccount string, gServiceAccount string) (err error) {
	const role = "roles/iam.workloadIdentityUser"
	var setIamPolicyBody []byte
	iamPolicy := iamPolicy{}
	binding := roleBinding{}

	binding.Role = role
	binding.Members = append(binding.Members, "serviceAccount:"+projectID+".svc.id.goog["+namespace+"/"+kServiceAccount+"]")

	iamPolicy.Bindings = append(iamPolicy.Bindings, binding)

	setIamPolicy := setIamPolicy{}
	setIamPolicy.Policy = iamPolicy
	if setIamPolicyBody, err = json.Marshal(setIamPolicy); err != nil {
		return err
	}

	u, _ := url.Parse(crmBetaURL)
	u.Path = path.Join(u.Path, GetProjectID()+":setIamPolicy")

	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())
	_, err = HttpClient(u.String(), string(setIamPolicyBody))

	return err
}
