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
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	"internal/apiclient"
	"internal/client/sites"
)

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

const maxPageSize = 100

type listapidocs struct {
	Status        string `json:"status,omitempty"`
	Message       string `json:"message,omitempty"`
	RequestID     string `json:"requestId,omitempty"`
	ErrorCode     string `json:"errorCode,omitempty"`
	Data          []data `json:"data,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

type apidocsdata struct {
	Status    string        `json:"status,omitempty"`
	Message   string        `json:"message,omitempty"`
	RequestID string        `json:"requestId,omitempty"`
	Data      documentation `json:"data,omitempty"`
}

type documentation struct {
	OasDocumentation     *oasDocumentation     `json:"oasDocumentation,omitempty"`
	GraphqlDocumentation *graphqlDocumentation `json:"graphqlDocumentation,omitempty"`
}

type oasDocumentation struct {
	Spec   spec   `json:"spec,omitempty"`
	Format string `json:"format,omitempty"`
}

type spec struct {
	DisplayName string `json:"displayName,omitempty"`
	Contents    string `json:"contents,omitempty"`
}

type graphqlDocumentation struct {
	Schema      spec   `json:"schema,omitempty"`
	EndpointUri string `json:"endpointUri,omitempty"`
}

type data struct {
	SiteID                   string   `json:"siteId,omitempty"`
	ID                       string   `json:"id,omitempty"`
	Title                    string   `json:"title,omitempty"`
	Description              string   `json:"description,omitempty"`
	Published                bool     `json:"published,omitempty"`
	AnonAllowed              bool     `json:"anonAllowed,omitempty"`
	ApiProductName           string   `json:"apiProductName,omitempty"`
	RequireCallbackUrl       bool     `json:"requireCallbackUrl,omitempty"`
	ImageUrl                 string   `json:"imageUrl,omitempty"`
	CategoryIDs              []string `json:"categoryIds,omitempty"`
	Modified                 string   `json:"modified,omitempty"`
	Visibility               bool     `json:"visibility,omitempty"`
	EdgeAPIProductName       string   `json:"edgeAPIProductName,omitempty"`
	SpecID                   string   `json:"specId,omitempty"`
	GraphqlSchema            string   `json:"graphqlSchema,omitempty"`
	GraphqlEndpointUrl       string   `json:"graphqlEndpointUrl,omitempty"`
	GraphqlSchemaDisplayName string   `json:"graphqlSchemaDisplayName,omitempty"`
}

type apidocResponse struct {
	Status    string `json:"status,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	ErrorCode string `json:"errorCode,omitempty"`
	Data      data   `json:"data,omitempty"`
}

// Create
func Create(siteid string, title string, description string, published string,
	anonAllowed string, apiProductName string, requireCallbackUrl string, imageUrl string,
	categoryIds []string,
) (respBody []byte, err error) {
	return createOrUpdate(siteid, "", title, description, published, anonAllowed,
		apiProductName, requireCallbackUrl, imageUrl, categoryIds, CREATE)
}

// Update
func Update(siteid string, id string, title string, description string, published string,
	anonAllowed string, apiProductName string, requireCallbackUrl string, imageUrl string,
	categoryIds []string,
) (respBody []byte, err error) {
	return createOrUpdate(siteid, id, title, description, published, anonAllowed,
		apiProductName, requireCallbackUrl, imageUrl, categoryIds, UPDATE)
}

func createOrUpdate(siteid string, id string, title string, description string, published string,
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
		apidoc = append(apidoc, "\"categoryIds\":[\""+getArrayStr(categoryIds)+"\"]")
	}

	payload := "{" + strings.Join(apidoc, ",") + "}"

	if action == CREATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs")
		respBody, err = apiclient.HttpClient(u.String(), payload)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", id)
		respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	}

	return respBody, err
}

// GetDocumentation
func GetDocumentation(siteid string, id string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", id, "documentation")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Get
func Get(siteid string, id string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", id)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// GetByTitle
func GetByTitle(siteid string, title string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	fullList := listapidocs{}
	pageToken := ""
	for {
		l := listapidocs{}
		listRespBytes, err := List(siteid, maxPageSize, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch apidocs: %w", err)
		}
		err = json.Unmarshal(listRespBytes, &l)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall: %w", err)
		}
		fullList.Data = append(fullList.Data, l.Data...)
		pageToken = l.NextPageToken
		if l.NextPageToken == "" {
			break
		}
	}
	for _, data := range fullList.Data {
		if data.Title == title {
			if respBody, err = json.Marshal(data); err != nil {
				return nil, err
			}
			return respBody, nil
		}
	}
	return nil, fmt.Errorf("unable to find apidocs with title %s", title)
}

// List
func List(siteid string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}

	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(siteid string, id string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", id)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// UpdateDocumentation
func UpdateDocumentation(siteid string, id string, displayName string,
	openAPIDoc string, graphQLDoc string, endpointUri string,
) (respBody []byte, err error) {
	var data map[string]interface{}
	// var payload string
	if openAPIDoc != "" {
		data = map[string]interface{}{
			"oasDocumentation": map[string]interface{}{
				"spec": map[string]interface{}{
					"displayName": displayName,
					"contents":    openAPIDoc,
				},
			},
		}
	}

	if graphQLDoc != "" {
		data = map[string]interface{}{
			"graphqlDocumentation": map[string]interface{}{
				"endpointUri": endpointUri,
				"schema": map[string]interface{}{
					"displayName": displayName,
					"contents":    graphQLDoc,
				},
			},
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}
	payload := string(jsonData)

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "sites", siteid, "apidocs", id, "documentation")
	respBody, err = apiclient.HttpClient(u.String(), payload, "PATCH")

	return nil, nil
}

// Export
func Export(folder string) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	siteids, err := sites.GetSiteIDs()
	if err != nil {
		return err
	}

	pageToken := ""

	for _, siteid := range siteids {
		l := listapidocs{}
		for {
			listRespBytes, err := List(siteid, maxPageSize, pageToken)
			if err != nil {
				return fmt.Errorf("failed to fetch apidocs: %w", err)
			}
			err = json.Unmarshal(listRespBytes, &l)
			if err != nil {
				return fmt.Errorf("failed to unmarshall: %w", err)
			}
			pageToken = l.NextPageToken
			// write apidocs Documentation
			for _, data := range l.Data {
				respDocsBody, err := GetDocumentation(siteid, data.ID)
				if err != nil {
					return err
				}
				docFileName := fmt.Sprintf("apidocs_%s_%s.json", siteid, data.ID)
				if err = apiclient.WriteByteArrayToFile(path.Join(folder, docFileName), false, respDocsBody); err != nil {
					return err
				}
			}
			if l.NextPageToken == "" {
				break
			}
		}
		respBody, err := json.Marshal(l.Data)
		if err != nil {
			return err
		}
		respBody, _ = apiclient.PrettifyJSON(respBody)
		if err = apiclient.WriteByteArrayToFile(path.Join(folder, "site_"+siteid+".json"), false, respBody); err != nil {
			return err
		}
	}
	return nil
}

/**
How the useNewSiteID flag works:
Assume data is exported from org1 which contains siteid site1. The fles are exported as
site_org1-site1.json and apidocs_org1-site1_00000.json

When importing this data into a new org, the siteid changes (since it is a combination of
org name and siteid).

Now import the data to org2 as
apigeecli apidocs import -o org2 -s site1 --source-f . --use-new-siteid=true -t $token
*/

// Import
func Import(siteid string, useSrcSiteID string, folder string) (err error) {
	var errs []string
	var respBody []byte
	var docsList []data

	if useSrcSiteID != "" {
		docsList, err = readAPIDocsDataFile(path.Join(folder, "site_"+useSrcSiteID+".json"))
	} else {
		docsList, err = readAPIDocsDataFile(path.Join(folder, "site_"+siteid+".json"))
	}
	if err != nil {
		return err
	}
	for _, doc := range docsList {
		// 1. create the apidoc object
		respBody, err = Create(siteid, doc.Title, doc.Description, strconv.FormatBool(doc.Published),
			strconv.FormatBool(doc.AnonAllowed), doc.ApiProductName,
			strconv.FormatBool(doc.RequireCallbackUrl), doc.ImageUrl, doc.CategoryIDs)
		if err != nil {
			errs = append(errs, err.Error())
			continue

		}

		// get the new doc.ID from the created apidoc
		response := apidocResponse{}
		err = json.Unmarshal(respBody, &response)
		if err != nil {
			return err
		}

		// 2. find the documentation associated with this site
		var documentationFileName string
		if useSrcSiteID != "" {
			apiDocsName := fmt.Sprintf("apidocs_%s_%s.json", useSrcSiteID, doc.ID)
			documentationFileName = path.Join(folder, apiDocsName)
		} else {
			documentationFileName = path.Join(folder, "apidocs_"+siteid+"_"+doc.ID+".json")
		}
		apidocument, err := readAPIDocumentationFile(documentationFileName)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if apidocument.Data.GraphqlDocumentation != nil {
			_, err = UpdateDocumentation(siteid, response.Data.ID,
				apidocument.Data.GraphqlDocumentation.Schema.DisplayName, "",
				apidocument.Data.GraphqlDocumentation.Schema.Contents,
				apidocument.Data.GraphqlDocumentation.EndpointUri)
		} else {
			_, err = UpdateDocumentation(siteid, response.Data.ID,
				apidocument.Data.OasDocumentation.Spec.DisplayName,
				apidocument.Data.OasDocumentation.Spec.Contents,
				"", "")
		}
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

func readAPIDocumentationFile(fileName string) (a apidocsdata, err error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return a, err
	}

	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return a, err
	}
	err = json.Unmarshal(content, &a)
	if err != nil {
		return a, err
	}
	return a, nil
}

func readAPIDocsDataFile(fileName string) (d []data, err error) {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		return d, err
	}

	defer jsonFile.Close()

	content, err := io.ReadAll(jsonFile)
	if err != nil {
		return d, err
	}
	err = json.Unmarshal(content, &d)
	if err != nil {
		return d, err
	}
	return d, nil
}
