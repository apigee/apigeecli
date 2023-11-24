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

package apidocs

import (
	"encoding/json"
	"net/url"
	"path"
	"strings"

	"internal/apiclient"
)

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

// Create
func Create(siteid string, title string, description string, published string,
	anonAllowed string, apiProductName string, requireCallbackUrl string, imageUrl string,
	categoryIds []string,
) (respBody []byte, err error) {
	return createOrUpdate(siteid, "", title, description, published, anonAllowed,
		apiProductName, requireCallbackUrl, imageUrl, categoryIds, CREATE)
}

// Update
func Update(siteid string, name string, title string, description string, published string,
	anonAllowed string, apiProductName string, requireCallbackUrl string, imageUrl string,
	categoryIds []string,
) (respBody []byte, err error) {
	return createOrUpdate(siteid, name, title, description, published, anonAllowed,
		apiProductName, requireCallbackUrl, imageUrl, categoryIds, UPDATE)
}

func createOrUpdate(siteid string, name string, title string, description string, published string,
	anonAllowed string, apiProductName string, requireCallbackUrl string, imageUrl string,
	categoryIds []string, action Action,
) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	apidoc := []string{}

	apidoc = append(apidoc, "\"siteId\":"+"\""+siteid+"\"")
	apidoc = append(apidoc, "\"title\":"+"\""+title+"\"")
	apidoc = append(apidoc, "\"apiProductName\":"+"\""+apiProductName+"\"")

	if description != "" {
		apidoc = append(apidoc, "\"description\":"+"\""+description+"\"")
	}

	if published != "" {
		apidoc = append(apidoc, "\"published\":"+published)
	}

	if anonAllowed != "" {
		apidoc = append(apidoc, "\"anonAllowed\":"+anonAllowed)
	}

	if requireCallbackUrl != "" {
		apidoc = append(apidoc, "\"requireCallbackUrl\":"+requireCallbackUrl)
	}

	if imageUrl != "" {
		apidoc = append(apidoc, "\"imageUrl\":"+"\""+imageUrl+"\"")
	}

	if len(categoryIds) > 0 {
		c, err := json.Marshal(categoryIds)
		if err != nil {
			return nil, err
		}
		categories := string(c)
		apidoc = append(apidoc, "\"categoryIds\":"+categories)
	}

	payload := "{" + strings.Join(apidoc, ",") + "}"

	if action == CREATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs")
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", name)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	}

	return respBody, err
}

// GetDocumentation
func GetDocumentation(siteid string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", name, "documentation")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Get
func Get(siteid string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(siteid string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(siteid string, name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// UpdateDocumentation
func UpdateDocumentation(siteid string, name string, openAPIDoc string, graphQLDoc string) (respBody []byte, err error) {
	var payload string

	if openAPIDoc != "" {
		payload = "{\"oasDocumentation\":\"spec\":{" + openAPIDoc + "}}"
	}

	if graphQLDoc != "" {
		payload = "{\"graphqlDocumentation\":\"spec\":{" + openAPIDoc + "}}"
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", name, "documentation")
	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")

	return nil, nil
}