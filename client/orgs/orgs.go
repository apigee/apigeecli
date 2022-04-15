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

package orgs

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/srinandan/apigeecli/apiclient"
)

var analyticsRegions = [...]string{"asia-east1", "asia-east1", "asia-northeast1", "asia-southeast1",
	"europe-west1", "us-central1", "us-east1", "us-east4", "us-west1", "australia-southeast1",
	"europe-west2"}

//OrgProperty contains an individual org flag or property
type orgProperty struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

//OrgProperties stores all the org feature flags and properties
type orgProperties struct {
	Property []orgProperty `json:"property,omitempty"`
}

type organization struct {
	Name                     string        `json:"name,omitempty"`
	DisplayName              string        `json:"displayName,omitempty"`
	Description              string        `json:"description,omitempty"`
	CreatedAt                string        `json:"createdAt,omitempty"`
	LastModifiedAt           string        `json:"lastModifiedAt,omitempty"`
	Environments             []string      `json:"environments,omitempty"`
	Properties               orgProperties `json:"properties,omitempty"`
	AnalyticsRegion          string        `json:"analyticsRegion,omitempty"`
	AuthorizedNetwork        string        `json:"authorizedNetwork,omitempty"`
	RuntimeType              string        `json:"runtimeType,omitempty"`
	SubscriptionType         string        `json:"subscriptionType,omitempty"`
	CaCertificate            string        `json:"caCertificate,omitempty"`
	RuntimeEncryptionKeyName string        `json:"runtimeDatabaseEncryptionKeyName,omitempty"`
	ProjectId                string        `json:"projectId,omitempty"`
	State                    string        `json:"state,omitempty"`
	BillingType              string        `json:"billingType,omitempty"`
}

func validRegion(region string) bool {
	for _, r := range analyticsRegions {
		if region == r {
			return true
		}
	}
	return false
}

//Create
func Create(region string, network string, runtimeType string, databaseKey string, billingType string, disablePortal bool) (respBody []byte, err error) {
	const baseURL = "https://apigee.googleapis.com/v1/organizations"
	var stageBaseURL = "https://staging-apigee.sandbox.googleapis.com/v1/organizations/"

	if !validRegion(region) {
		return respBody, fmt.Errorf("invalid analytics region."+
			" Analytics region must be one of : %v", analyticsRegions)
	}

	var u *url.URL
	if apiclient.GetStaging() {
		u, _ = url.Parse(stageBaseURL)
	} else {
		u, _ = url.Parse(baseURL)
	}

	u.Path = path.Join(u.Path)
	q := u.Query()
	q.Set("parent", "projects/"+apiclient.GetProjectID())
	u.RawQuery = q.Encode()

	orgPayload := []string{}
	orgPayload = append(orgPayload, "\"name\":\""+apiclient.GetApigeeOrg()+"\"")
	orgPayload = append(orgPayload, "\"analyticsRegion\":\""+region+"\"")
	orgPayload = append(orgPayload, "\"runtimeType\":\""+runtimeType+"\"")
	if disablePortal {
		orgPayload = append(orgPayload, "\"portalDisabled\": true")
	}
	if runtimeType == "CLOUD" {
		orgPayload = append(orgPayload, "\"authorizedNetwork\":\""+network+"\"")
		orgPayload = append(orgPayload, "\"runtimeDatabaseEncryptionKeyName\":\""+databaseKey+"\"")
	}

	if billingType != "" {
		orgPayload = append(orgPayload, "\"billingType\":\""+billingType+"\"")
	}

	payload := "{" + strings.Join(orgPayload, ",") + "}"
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)
	return respBody, err
}

//Get
func Get() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg())
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

func GetOrgField(key string) (value string, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg())

	apiclient.SetPrintOutput(false)
	orgBody, err := apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	if err != nil {
		return "", err
	}
	apiclient.SetPrintOutput(true)

	var orgMap map[string]interface{}
	err = json.Unmarshal(orgBody, &orgMap)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", orgMap[key]), nil
}

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//GetDeployedIngressConfig
func GetDeployedIngressConfig(view bool) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if view {
		q := u.Query()
		q.Set("view", "full")
		u.RawQuery = q.Encode()
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "deployedIngressConfig")
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//SetOrgProperty is used to set org properties
func SetOrgProperty(name string, value string) (err error) {

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg())
	//get org details
	orgBody, err := apiclient.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	org := organization{}
	err = json.Unmarshal(orgBody, &org)
	if err != nil {
		return err
	}

	//check if the property exists
	found := false
	for i, properties := range org.Properties.Property {
		if properties.Name == name {
			fmt.Println("Property found, enabling property")
			org.Properties.Property[i].Value = value
			found = true
			break
		}
	}

	if !found {
		//set the property
		newProp := orgProperty{}
		newProp.Name = name
		newProp.Value = value

		org.Properties.Property = append(org.Properties.Property, newProp)
	}

	newOrgBody, err := json.Marshal(org)
	if err != nil {
		return err
	}

	u, _ = url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg())
	_, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(newOrgBody), "PUT")

	return err
}

//Update
func Update(description string, displayName string, region string, network string, runtimeType string, databaseKey string) (respBody []byte, err error) {

	apiclient.SetPrintOutput(false)
	orgBody, err := Get()
	if err != nil {
		return nil, err
	}
	apiclient.SetPrintOutput(true)

	org := organization{}
	err = json.Unmarshal(orgBody, &org)
	if err != nil {
		return nil, err
	}

	if description != "" {
		org.Description = description
	}

	if displayName != "" {
		org.DisplayName = displayName
	}

	if region != "" {
		org.AnalyticsRegion = region
	}

	if network != "" {
		org.AuthorizedNetwork = network
	}

	if runtimeType != "" {
		org.RuntimeType = runtimeType
	}

	if databaseKey != "" {
		org.RuntimeEncryptionKeyName = databaseKey
	}

	newOrgBody, err := json.Marshal(org)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg())
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), string(newOrgBody), "PUT")

	return respBody, err
}

//SetAddons
func SetAddons(advancedApiOpsConfig bool, integrationConfig bool, monetizationConfig bool, connectorsConfig bool) (respBody []byte, err error) {

	addonPayload := []string{}

	if !advancedApiOpsConfig && !integrationConfig && !monetizationConfig {
		return nil, fmt.Errorf("At least one addon must be enabled")
	}

	if advancedApiOpsConfig {
		addonPayload = append(addonPayload, "\"advancedApiOpsConfig\":{\"enabled\":true}")
	}

	if integrationConfig {
		addonPayload = append(addonPayload, "\"integrationConfig\":{\"enabled\":true}")
	}

	if monetizationConfig {
		addonPayload = append(addonPayload, "\"monetizationConfig\":{\"enabled\":true}")
	}

	if connectorsConfig {
		addonPayload = append(addonPayload, "\"connectorsPlatformConfig\":{\"enabled\":true}")
	}

	payload := "{\"addonsConfig\":{" + strings.Join(addonPayload, ",") + "}}"

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg()+":setAddons")

	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String(), payload)

	return respBody, err
}
