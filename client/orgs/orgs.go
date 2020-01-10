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

func validRegion(region string) bool {
	for _, r := range analyticsRegions {
		if region == r {
			return true
		}
	}
	return false
}

//Create
func Create(region string) (respBody []byte, err error) {
	const baseURL = "https://apigee.googleapis.com/v1/organizations"

	if !validRegion(region) {
		return respBody, fmt.Errorf("invalid analytics region."+
			" Analytics region must be one of : %v", analyticsRegions)
	}

	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path)
	q := u.Query()
	q.Set("parent", "projects/"+apiclient.GetProjectID())
	u.RawQuery = q.Encode()

	orgPayload := []string{}
	orgPayload = append(orgPayload, "\"name\":\""+apiclient.GetApigeeOrg()+"\"")
	orgPayload = append(orgPayload, "\"analyticsRegion\":\""+region+"\"")

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

//List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	respBody, err = apiclient.HttpClient(apiclient.GetPrintOutput(), u.String())
	return respBody, err
}

//SetOrgProperty is used to set org properties
func SetOrgProperty(name string, value string) (err error) {
	//OrgProperty contains an individual org flag or property
	type orgProperty struct {
		Name  string `json:"name,omitempty"`
		Value string `json:"value,omitempty"`
	}
	//OrgProperties stores all the org feature flags and properties
	type orgProperties struct {
		Property []orgProperty `json:"property,omitempty"`
	}
	//Org structure
	type organization struct {
		Name            string        `json:"name,omitempty"`
		CreatedAt       string        `json:"-,omitempty"`
		LastModifiedAt  string        `json:"-,omitempty"`
		Environments    []string      `json:"-,omitempty"`
		Properties      orgProperties `json:"properties,omitempty"`
		AnalyticsRegion string        `json:"-,omitempty"`
	}

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
