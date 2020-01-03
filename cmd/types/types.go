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

package types

//Attribute to used to hold custom attributes for entities
type Attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Arguments is the base struct to hold all command arguments
type Arguments struct {
	Org            string //Apigee org
	Env            string //Apigee environment
	Token          string //Google OAuth access token
	ServiceAccount string //Google service account json
	AliasName      string //AliasName for the key store
	ProjectID      string //GCP Project ID
	LogInfo        bool   //LogInfo controls the log level
	SkipCheck      bool   //skip checking access token expiry
	SkipCache      bool   //skip writing access token to file
}

//OAuthAccessToken is a structure to hold OAuth response
type OAuthAccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

//KeyAliasName holds the name of the key alias
type KeyAliasName string

//Binding for IAM Roles
type Binding struct {
	Role      string     `json:"role,omitempty"`
	Members   []string   `json:"members,omitempty"`
	Condition *Condition `json:"condition,omitempty"`
}

//Condition for Bindings
type Condition struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Expression  string `json:"expression,omitempty"`
}

//IamPolicy
type IamPolicy struct {
	Version  int       `json:"version,omitempty"`
	Etag     string    `json:"etag,omitempty"`
	Bindings []Binding `json:"bindings,omitempty"`
}

//SetIamPolicy
type SetIamPolicy struct {
	Policy IamPolicy `json:"policy,omitempty"`
}

//Org structure
type Org struct {
	Name            string        `json:"name,omitempty"`
	CreatedAt       string        `json:"-,omitempty"`
	LastModifiedAt  string        `json:"-,omitempty"`
	Environments    []string      `json:"-,omitempty"`
	Properties      OrgProperties `json:"properties,omitempty"`
	AnalyticsRegion string        `json:"-,omitempty"`
}

//OrgProperties stores all the org feature flags and properties
type OrgProperties struct {
	Property []OrgProperty `json:"property,omitempty"`
}

//OrgProperty contains an individual org flag or property
type OrgProperty struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

func (a KeyAliasName) String() string {
	return string(a)
}

//ResourceTypes contains a list of valid resources
var resourceTypes = [7]string{"js", "jsc", "properties", "java", "wsdl", "xsd", "py"}

//IsValidResource returns true is the resource type is valid
func IsValidResource(resType string) bool {
	for _, n := range resourceTypes {
		if n == resType {

			return true
		}
	}

	return false
}
