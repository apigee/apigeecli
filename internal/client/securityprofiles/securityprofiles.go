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
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
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

type computeenvscore struct {
	TimeRange timeInterval  `json:"timeRange,omitempty"`
	Filters   []scoreFilter `json:"filters,omitempty"`
	PageSize  int           `json:"pageSize,omitempty"`
	PageToken string        `json:"pageToken,omitempty"`
}

type timeInterval struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}

type scoreFilter struct {
	ScorePath string `json:"scorePath,omitempty"`
}

type scores struct {
	Scores        []score `json:"scores,omitempty"`
	NextPageToken string  `json:"nextPageToken,omitempty"`
}

type score struct {
	TimeRange     timeInterval `json:"timeRange,omitempty"`
	Component     component    `json:"component,omitempty"`
	Subcomponents []component  `json:"subcomponents,omitempty"`
}

type component struct {
	Score           int              `json:"score,omitempty"`
	ScorePath       string           `json:"scorePath,omitempty"`
	Recommendations []recommendation `json:"recommendations,omitempty"`
	DataCaptureTime string           `json:"dataCaptureTime,omitempty"`
	CalculateTime   string           `json:"calculateTime,omitempty"`
	DrilldownPaths  []string         `json:"drilldownPaths,omitempty"`
}

type recommendation struct {
	Title       string   `json:"title,omitempty"`
	Description string   `json:"description,omitempty"`
	Impact      int      `json:"impact,omitempty"`
	Actions     []action `json:"actions,omitempty"`
}

type action struct {
	Description   string        `json:"description,omitempty"`
	ActionContext actionContext `json:"actionContext,omitempty"`
}

type actionContext struct {
	DocumentationLink string `json:"documentationLink,omitempty"`
}

const maxPageSize = 50

// Create
func Create(name string, content []byte) (respBody []byte, err error) {
	sc := secprofile{}
	if err = json.Unmarshal(content, &sc); err != nil {
		return nil, err
	}
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles")
	q := u.Query()
	q.Set("securityProfileId", name)
	u.RawQuery = q.Encode()
	respBody, err = apiclient.HttpClient(u.String(), string(content))
	return respBody, err
}

// Attach
func Attach(name string, revision string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles",
		name, "environments", apiclient.GetApigeeEnv())
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// Get
func Get(name string, revision string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name+":listRevisions")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// List
func List(pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
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
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles", name)
	q := u.Query()
	q.Set("updateMask", "description,profileConfig")
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), string(content), "PATCH")
	return respBody, err
}

// Compute
func Compute(name string, startTime string, endTime string, filters []string,
	pageSize int, pageToken string, full bool,
) (respBody []byte, err error) {
	score := computeenvscore{}
	if pageSize != -1 {
		score.PageSize = pageSize
	}
	if pageToken != "" {
		score.PageToken = pageToken
	}
	if startTime == "" || endTime == "" {
		return nil, fmt.Errorf("startTime and endTime are mandatory")
	}
	score.TimeRange.StartTime = startTime
	score.TimeRange.EndTime = endTime

	for _, f := range filters {
		sf := scoreFilter{}
		sf.ScorePath = f
		score.Filters = append(score.Filters, sf)
	}

	fullScores := scores{}

	for {
		payload, err := json.Marshal(score)
		if err != nil {
			return nil, err
		}

		u, _ := url.Parse(apiclient.GetApigeeBaseURL())
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "securityProfiles",
			name, "environments", apiclient.GetApigeeEnv()+":computeEnvironmentScores")

		if respBody, err = apiclient.HttpClient(u.String(), string(payload)); err != nil {
			return nil, err
		}

		if !full {
			return respBody, err
		}

		s := scores{}
		if err = json.Unmarshal(respBody, &s); err != nil {
			return nil, err
		}

		fullScores.Scores = append(fullScores.Scores, s.Scores...)
		score.PageToken = s.NextPageToken

		if s.NextPageToken == "" {
			break
		}

	}

	respBody, err = json.Marshal(fullScores)
	return respBody, err
}

// Export
func Export(conn int, folder string) (err error) {
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
		pageToken = l.NextPageToken
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
		go exportWorker(&fanOutWg, workChan, folder, errChan)
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

func exportWorker(wg *sync.WaitGroup, workCh <-chan secprofile, folder string, errs chan<- error) {
	defer wg.Done()
	for {
		var err error

		work, ok := <-workCh
		if !ok {
			return
		}

		if work.Name == "default" {
			// do not export the default profile. this is set by the system
			return
		}

		clilog.Info.Printf("Exporting Security Profile %s\n", work.Name)
		payload, err := json.Marshal(work)
		if err != nil {
			errs <- err
		}
		payload, _ = apiclient.PrettifyJSON(payload)
		apiclient.WriteByteArrayToFile(path.Join(folder, work.Name+".json"), false, payload)
	}
}

// Import
func Import(conn int, folder string) error {
	var profiles []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".json" {
			return nil
		}
		profiles = append(profiles, path)
		return nil
	})
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		clilog.Warning.Println("No profiles were found in the folder")
		return nil
	}

	clilog.Debug.Printf("Found %d profiles in the folder\n", len(profiles))
	clilog.Debug.Printf("Importing security profiles with %d parallel connections\n", conn)

	jobChan := make(chan string)
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
		go importProfile(&fanOutWg, jobChan, errChan)
	}

	for _, profile := range profiles {
		jobChan <- profile
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

func importProfile(wg *sync.WaitGroup, jobs <-chan string, errs chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		content, err := readFile(job)
		if err != nil {
			errs <- err
			continue
		}
		if _, err = Create(getSecurityProfileName(job), content); err != nil {
			errs <- err
			continue
		}
	}
}

func readFile(filePath string) (byteValue []byte, err error) {
	userFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer userFile.Close()

	byteValue, err = io.ReadAll(userFile)
	if err != nil {
		return nil, err
	}
	return byteValue, err
}

func getSecurityProfileName(job string) (name string) {
	ver := regexp.MustCompile(`@\d+`)
	return ver.ReplaceAllString(strings.TrimSuffix(filepath.Base(job), filepath.Ext(job)), "")
}
