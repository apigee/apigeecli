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

package securityprofiles

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"
	"internal/clilog"
)

type secprofiles struct {
	SecurityProfiles []secprofile `json:"securityProfiles,omitempty"`
	NextPageToken    string       `json:"nextPageToken,omitempty"`
}

type secprofile struct {
	Name          string        `json:"name"`
	DisplayName   string        `json:"displayName"`
	Description   string        `json:"description,omitempty"`
	RevisionID    string        `json:"revisionId,omitempty"`
	ProfileConfig profileConfig `json:"profileConfig"`
	ScoreConfigs  []scoreConfig `json:"scoringConfigs,omitempty"`
}

type profileConfig struct {
	Categories []category `json:"categories"`
}

type scoreConfig struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ScorePath   string `json:"scorePath,omitempty"`
}

type category struct {
	Abuse         interface{} `json:"abuse,omitempty"`
	Mediation     interface{} `json:"mediation,omitempty"`
	Authorization interface{} `json:"authorization,omitempty"`
	Threat        interface{} `json:"threat,omitempty"`
	Mtls          interface{} `json:"mtls,omitempty"`
	Cors          interface{} `json:"cors,omitempty"`
}

const maxPageSize = 50

// Create
func Create(name string, content []byte) (respBody []byte, err error) {
	sc := secprofile{}
	if err = json.Unmarshal(content, &sc); err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles")
	q := u.Query()
	q.Set("securityProfileId", name)
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), string(content))
	return respBody, err
}

// Attach
func Attach(name string, revision string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name, "environments")
	attachProfile := []string{}
	attachProfile = append(attachProfile, "\"name\":"+"\""+apiclient.GetApigeeEnv()+"\"")
	attachProfile = append(attachProfile, "\"securityProfileRevisionId\":"+"\""+revision+"\"")
	payload := "{" + strings.Join(attachProfile, ",") + "}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Detach
func Detach(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles",
		name, "environments", apiclient.GetApigeeEnv())
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(name string, revision string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	if revision != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name+"@"+revision)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListRevisions
func ListRevisions(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name+":listRevisions")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles")
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

// Update
func Update(name string, content []byte) (respBody []byte, err error) {
	sc := secprofile{}
	if err = json.Unmarshal(content, &sc); err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	q := u.Query()
	q.Set("updateMask", "description,profileConfig")
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), string(content), "PATCH")
	return respBody, err
}

// Export
func Export(conn int, folder string, allRevisions bool) (err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	pageToken := ""
	listsecprofiles := secprofiles{}

	for {
		l := secprofiles{}
		listRespBytes, err := List(maxPageSize, pageToken)
		if err != nil {
			return fmt.Errorf("failed to fetch security profiles: %w", err)
		}
		err = json.Unmarshal(listRespBytes, &l)
		if err != nil {
			return fmt.Errorf("failed to unmarshall: %w", err)
		}
		listsecprofiles.SecurityProfiles = append(listsecprofiles.SecurityProfiles, l.SecurityProfiles...)
		if l.NextPageToken == "" {
			break
		}
	}

	errChan := make(chan error)
	workChan := make(chan secprofile, len(listsecprofiles.SecurityProfiles))

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
		go exportWorker(&fanOutWg, workChan, folder, allRevisions, errChan)
	}

	for _, i := range listsecprofiles.SecurityProfiles {
		workChan <- i
	}

	close(workChan)
	fanOutWg.Wait()
	close(errChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}

	return nil
}

func exportWorker(wg *sync.WaitGroup, workCh <-chan secprofile, folder string, allRevisions bool, errs chan<- error) {
	defer wg.Done()
	for {
		var respBody []byte
		var err error

		listsecprofiles := secprofiles{}

		work, ok := <-workCh
		if !ok {
			return
		}

		if allRevisions {
			securityProfileName := work.Name[strings.LastIndex(work.Name, "/")+1:]
			clilog.Info.Printf("Exporting all the revisions for Security Profile %s\n", securityProfileName)

			if respBody, err = ListRevisions(securityProfileName); err != nil {
				errs <- err
			}
			err = json.Unmarshal(respBody, &listsecprofiles)
			if err != nil {
				errs <- err
			}
			for _, s := range listsecprofiles.SecurityProfiles {
				payload, err := json.Marshal(s)
				if err != nil {
					errs <- err
				}
				payload, _ = apiclient.PrettifyJSON(payload)
				apiclient.WriteByteArrayToFile(path.Join(folder, s.Name), false, payload)
			}
		} else {
			clilog.Info.Printf("Exporting Security Profile %s\n", work.Name)
			payload, err := json.Marshal(work)
			if err != nil {
				errs <- err
			}
			payload, _ = apiclient.PrettifyJSON(payload)
			apiclient.WriteByteArrayToFile(path.Join(folder, work.Name+".json"), false, payload)
		}
	}
}
