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

package appgroups

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
	"sync"

	"internal/apiclient"
	"internal/clilog"
)

type appgroups struct {
	AppGroups     []appgroup `json:"appGroups,omitempty"`
	NextPageToken string     `json:"nextPageToken,omitempty"`
}

type appgroup struct {
	Name        string      `json:"name,omitempty"`
	Status      *string     `json:"status,omitempty"`
	AppGroupId  string      `json:"appGroupId,omitempty"`
	ChannelURI  string      `json:"channelUri,omitempty"`
	ChannelID   string      `json:"channelId,omitempty"`
	Attributes  []attribute `json:"attributes,omitempty"`
	DisplayName string      `json:"displayName,omitempty"`
}

// attribute to used to hold custom attributes for entities
type attribute struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type devDetail struct {
	DeveloperId string   `json:"developerId,omitempty"`
	Roles       []string `json:"roles,omitempty"`
}

var maxPageSize = 1000

// Create
func Create(name string, channelURI string, channelID string, displayName string, attrs map[string]string, devs map[string]string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	app := []string{}

	app = append(app, "\"name\":\""+name+"\"")
	if channelURI != "" {
		app = append(app, "\"channelUri\":\""+channelURI+"\"")
	}

	if channelID != "" {
		app = append(app, "\"channelId\":\""+channelID+"\"")
	}

	if displayName != "" {
		app = append(app, "\"displayName\":\""+displayName+"\"")
	}

	if len(devs) > 0 {
		if len(attrs) == 0 {
			attrs = make(map[string]string)
		}
		attrs["__apigee_reserved__developer_details"] = strings.ReplaceAll(getDeveloperDetails(devs), "\"", "\\\"")
	}

	if len(attrs) != 0 {
		attributes := []string{}
		for key, value := range attrs {
			attributes = append(attributes, "{\"name\":\""+key+"\",\"value\":\""+value+"\"}")
		}
		attributesStr := "\"attributes\":[" + strings.Join(attributes, ",") + "]"
		app = append(app, attributesStr)
	}

	payload := "{" + strings.Join(app, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Update
func Update(name string, channelURI string, channelID string, displayName string, attrs map[string]string, devs map[string]string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	appGroupRespBody, err := Get(name)
	developerDetailsAttribute := attribute{}

	if err != nil {
		return nil, err
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	a := appgroup{}
	if err = json.Unmarshal(appGroupRespBody, &a); err != nil {
		return nil, err
	}

	if channelURI != "" {
		a.ChannelURI = channelURI
	}

	if channelID != "" {
		a.ChannelID = channelID
	}

	if displayName != "" {
		a.DisplayName = displayName
	}

	if len(devs) > 0 {
		if len(attrs) == 0 {
			attrs = make(map[string]string)
		}
		developerDetailsAttribute.Name = "__apigee_reserved__developer_details"
		developerDetailsAttribute.Value = getDeveloperDetails(devs)
	}

	// clear existing attributes
	a.Attributes = nil

	if len(attrs) != 0 {
		for key, value := range attrs {
			attr := attribute{}
			attr.Name = key
			attr.Value = value
			a.Attributes = append(a.Attributes, attr)
		}
	}

	if developerDetailsAttribute != (attribute{}) {
		a.Attributes = append(a.Attributes, developerDetailsAttribute)
	}

	reqBody, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	respBody, err = apiclient.HttpClient(u.String(), string(reqBody), "PUT")
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Manage
func Manage(name string, action string) (respBody []byte, err error) {
	if action != "active" && action != "inactive" {
		return nil, fmt.Errorf("invalid action. action must be active or inactive")
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name)
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), "", "PUT", "application/octet-stream")
	return respBody, err
}

// List
func List(pageSize int, pageToken string, filter string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups")
	q := u.Query()
	if pageSize != -1 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}
	if pageToken != "" {
		q.Set("pageToken", pageToken)
	}
	if filter != "" {
		q.Set("filter", filter)
	}
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Export
func Export() (respBody []byte, err error) {
	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	pageToken := ""
	applist := appgroups{}

	for {
		a := appgroups{}
		listRespBytes, err := List(maxPageSize, pageToken, "")
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(listRespBytes, &a)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall: %w", err)
		}
		applist.AppGroups = append(applist.AppGroups, a.AppGroups...)
		pageToken = a.NextPageToken
		if a.NextPageToken == "" {
			break
		}
	}

	respBody, err = json.Marshal(applist.AppGroups)

	return respBody, err
}

// Import
func Import(conn int, filePath string) error {
	appgrps, err := readAppGroupsFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Debug.Printf("Found %d AppGroups in the file\n", len(appgrps))
	clilog.Debug.Printf("Create AppGroups with %d connections\n", conn)

	knownAppGroupsBytes, err := Export()
	if err != nil {
		return err
	}

	var knownAppGroups []appgroup
	if err = json.Unmarshal(knownAppGroupsBytes, &knownAppGroups); err != nil {
		return err
	}

	knownAppGroupsList := map[string]bool{}
	for _, name := range knownAppGroups {
		knownAppGroupsList[name.Name] = true
	}

	jobChan := make(chan appgroup)
	errChan := make(chan error)

	fanOutWg := sync.WaitGroup{}
	fanInWg := sync.WaitGroup{}

	errs := []string{}
	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newErr, ok := <-errChan
			if !ok {
				return
			}
			errs = append(errs, newErr.Error())
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go importAppGroup(knownAppGroupsList, &fanOutWg, jobChan, errChan)
	}

	for _, ag := range appgrps {
		jobChan <- ag
	}

	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func importAppGroup(knownAppGroupsList map[string]bool, wg *sync.WaitGroup, jobs <-chan appgroup, errs chan<- error) {
	defer wg.Done()
	var err error

	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		if knownAppGroupsList[job.Name] {
			// the appgroup already exists, perform an update
			_, err = Update(job.Name, job.ChannelURI, job.ChannelID, job.DisplayName, getMapAttributes(job.Attributes), nil)
		} else {
			_, err = Create(job.Name, job.ChannelURI, job.ChannelID, job.DisplayName, getMapAttributes(job.Attributes), nil)
		}

		if err != nil {
			errs <- err
			continue
		}

		clilog.Debug.Printf("Completed appgroup: %s", job.Name)
	}
}

func readAppGroupsFile(filePath string) ([]appgroup, error) {
	a := []appgroup{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return a, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return a, err
	}

	err = json.Unmarshal(byteValue, &a)
	if err != nil {
		return a, err
	}

	return a, nil
}

func getMapAttributes(attrs []attribute) map[string]string {
	var custAttrs map[string]string
	for _, a := range attrs {
		custAttrs[a.Name] = a.Value
	}
	return custAttrs
}

func getDeveloperDetails(devs map[string]string) string {
	developerDetails := []string{}
	for k, v := range devs {
		developerDetails = append(developerDetails, "{ \"developerId\":\""+k+"\", \"roles\":[ \""+v+"\" ] }")
	}
	// The docs https://cloud.google.com/apigee/docs/reference/apis/apigee/rest/v1/organizations.appgroups/create
	// called for base64 encoding. However, it appears the Apigee Durpal implementation does not require it.
	// return base64.StdEncoding.EncodeToString([]byte("[" + strings.Join(developerDetails, ",") + "]"))
	developerDetailsString := "[" + strings.Join(developerDetails, ",") + "]"
	return developerDetailsString
}

func convertToAttributes(a map[string]string) []attribute {
	attributes := []attribute{}
	for k, v := range a {
		attr := attribute{}
		attr.Name = k
		attr.Value = v
		attributes = append(attributes, attr)
	}
	return attributes
}
