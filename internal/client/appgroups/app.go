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

type appgroupapps struct {
	AppGroupApps  []appgroupapp `json:"appGroupApps,omitempty"`
	NextPageToken string        `json:"nextPageToken,omitempty"`
}

type appgroupapp struct {
	AppId          string       `json:"appId,omitempty"`
	Name           string       `json:"name,omitempty"`
	Status         string       `json:"status,omitempty"`
	AppGroup       string       `json:"appGroup,omitempty"`
	CreatedAt      string       `json:"createdAt,omitempty"`
	LastModifiedAt string       `json:"lastModifiedAt,omitempty"`
	Credentials    []credential `json:"credentials,omitempty"`
}

type credential struct {
	APIProducts    []apiProduct `json:"apiProducts,omitempty"`
	ConsumerKey    string       `json:"consumerKey,omitempty"`
	ConsumerSecret string       `json:"consumerSecret,omitempty"`
	ExpiresAt      string       `json:"expiresAt,omitempty"`
	Status         string       `json:"status,omitempty"`
	Scopes         []string     `json:"scopes,omitempty"`
}

type apiProduct struct {
	ApiProduct string `json:"apiproduct,omitempty"`
	Status     string `json:"status,omitempty"`
}

type Action uint8

const (
	CREATE Action = iota
	UPDATE
)

// CreateApp
func CreateApp(name string, appName string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createOrUpdate(name, appName, expires, callback, apiProducts, scopes, attrs, CREATE)
}

// createAppNoKey
func createAppNoKey(name string, appName string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (err error) {
	respBody, err := CreateApp(name, appName, expires, callback, apiProducts, scopes, attrs)
	if err != nil {
		return err
	}

	a := appgroupapp{}
	err = json.Unmarshal(respBody, &a)
	if err != nil {
		return err
	}

	_, err = DeleteKey(name, appName, a.Credentials[0].ConsumerKey)
	return err
}

// DeleteApp
func DeleteApp(name string, appName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetApp
func GetApp(name string, appName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ListApps
func ListApps(name string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps")
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

// UpdateApp
func UpdateApp(name string, appName string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string) (respBody []byte, err error) {
	return createOrUpdate(name, appName, expires, callback, apiProducts, scopes, attrs, UPDATE)
}

func createOrUpdate(name string, appName string, expires string, callback string, apiProducts []string, scopes []string, attrs map[string]string, action Action) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	app := []string{}

	app = append(app, "\"name\":\""+appName+"\"")

	if len(apiProducts) > 0 {
		app = append(app, "\"apiProducts\":[\""+getArrayStr(apiProducts)+"\"]")
	}

	if callback != "" {
		app = append(app, "\"callbackUrl\":\""+callback+"\"")
	}

	if expires != "" {
		app = append(app, "\"keyExpiresIn\":\""+expires+"\"")
	}

	if len(scopes) > 0 {
		app = append(app, "\"scopes\":[\""+getArrayStr(scopes)+"\"]")
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

	if action == CREATE {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps")
		respBody, err = apiclient.HttpClient(u.String(), payload)
		return respBody, err
	}
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	return respBody, err
}

// ManageApp
func ManageApp(name string, appName string, action string) (respBody []byte, err error) {
	if action != "revoke" && action != "approve" {
		return nil, fmt.Errorf("invalid action. action must be revoke or approve")
	}

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "appgroups", name, "apps", appName)
	q := u.Query()
	q.Set("action", action)
	u.RawQuery = q.Encode()

	respBody, err = apiclient.HttpClient(u.String(), "", "PUT", "application/octet-stream")
	return respBody, err
}

// ExportApps
func ExportApps(name string) (respBody []byte, err error) {
	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	pageToken := ""
	applist := appgroupapps{}

	for {
		a := appgroupapps{}
		listRespBytes, err := ListApps(name, maxPageSize, pageToken)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(listRespBytes, &a)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshall: %w", err)
		}
		applist.AppGroupApps = append(applist.AppGroupApps, a.AppGroupApps...)
		pageToken = a.NextPageToken
		if a.NextPageToken == "" {
			break
		}
	}

	respBody, err = json.Marshal(applist.AppGroupApps)

	return respBody, err
}

// ImportApps
func ImportApps(conn int, filePath string, name string) (err error) {
	appgrpapps, err := readAppsFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Debug.Printf("Found %d apps in the file\n", len(appgrpapps))
	clilog.Debug.Printf("Create apps with %d connections\n", conn)

	knownAppGroupAppsBytes, err := ExportApps(name)
	if err != nil {
		return err
	}

	var knownAppGroupApps []appgroupapp
	if err = json.Unmarshal(knownAppGroupAppsBytes, &knownAppGroupApps); err != nil {
		return err
	}

	knownAppGroupAppsList := map[string]bool{}
	for _, name := range knownAppGroupApps {
		knownAppGroupAppsList[name.Name] = true
	}

	jobChan := make(chan appgroupapp)
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
		go importAppGroupApp(knownAppGroupAppsList, &fanOutWg, jobChan, name, errChan)
	}

	for _, aga := range appgrpapps {
		jobChan <- aga
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

func importAppGroupApp(knownAppGroupsList map[string]bool, wg *sync.WaitGroup, jobs <-chan appgroupapp, name string, errs chan<- error) {
	defer wg.Done()
	var err error

	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		if knownAppGroupsList[job.Name] {
			// the appgroup already exists, perform an update
			clilog.Warning.Printf("App %s in AppGroup %s already exists. Updates to app is not currently supported", job.Name, name)
		} else {
			// 1. Create the app
			err = createAppNoKey(name, job.Name, "-1", "", nil, nil, nil)
			if err != nil {
				errs <- err
				continue
			}
			// 2. Create app keys
			for _, c := range job.Credentials {
				_, err = CreateKey(name, job.Name, c.ConsumerKey, c.ConsumerSecret, c.ExpiresAt, getAPIProductsAsArray(c.APIProducts), c.Scopes, nil)
				if err != nil {
					errs <- err
				}
			}
		}

		clilog.Debug.Printf("Completed app %s import to appgroup %s", job.Name, name)
	}
}

// ExportAllApps
func ExportAllApps(appGroupListBytes []byte, conn int) (results [][]byte, err error) {
	appGroupsList := []appgroup{}

	err = json.Unmarshal(appGroupListBytes, &appGroupsList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall: %w", err)
	}

	clilog.Debug.Printf("Found %d appgroups in the org\n", len(appGroupsList))
	clilog.Debug.Printf("Exporting appgroups with %d parallel connections\n", conn)

	jobChan := make(chan appgroup)
	resultChan := make(chan []byte)
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

	fanInWg.Add(1)
	go func() {
		defer fanInWg.Done()
		for {
			newResult, ok := <-resultChan
			if !ok {
				return
			}
			results = append(results, newResult)
		}
	}()

	for i := 0; i < conn; i++ {
		fanOutWg.Add(1)
		go exportApp(&fanOutWg, jobChan, resultChan, errChan)
	}

	for _, a := range appGroupsList {
		jobChan <- a
	}
	close(jobChan)
	fanOutWg.Wait()
	close(errChan)
	close(resultChan)
	fanInWg.Wait()

	if len(errs) > 0 {
		return nil, errors.New(strings.Join(errs, "\n"))
	}
	return results, nil
}

// ImportAllApps
func ImportAllApps(conn int, filePath string) (err error) {
	appgrpapps, err := readAppsArrayFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Info.Printf("Found %d apps in the file\n", len(appgrpapps))
	clilog.Info.Printf("Create apps with %d connections\n", conn)

	// flatten the structure
	a := []appgroupapp{}
	for _, aga := range appgrpapps {
		a = append(a, aga[0])
	}

	jobChan := make(chan appgroupapp)
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
		go importAllAppGroupApps(&fanOutWg, jobChan, errChan)
	}

	for _, aga := range a {
		jobChan <- aga
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

func importAllAppGroupApps(wg *sync.WaitGroup, jobs <-chan appgroupapp, errs chan<- error) {
	defer wg.Done()
	var err error

	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		// 1. Create the app
		err = createAppNoKey(job.AppGroup, job.Name, "-1", "", nil, nil, nil)
		if err != nil {
			errs <- err
			continue
		}
		// 2. Create app keys
		for _, c := range job.Credentials {
			_, err = CreateKey(job.AppGroup, job.Name, c.ConsumerKey, c.ConsumerSecret, c.ExpiresAt, getAPIProductsAsArray(c.APIProducts), c.Scopes, nil)
			if err != nil {
				errs <- err
			}
		}
		clilog.Debug.Printf("Completed app %s import to appgroup %s", job.Name, job.AppGroup)
	}
}

func exportApp(wg *sync.WaitGroup, jobs <-chan appgroup, results chan<- []byte, errs chan<- error) {
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		respBody, err := ExportApps(job.Name)
		if err != nil {
			errs <- err
		} else {
			results <- respBody
			clilog.Debug.Printf("Completed Exporting app %s in appgroup %s\n", job.Name, job.AppGroupId)
		}
	}
}

func getArrayStr(str []string) string {
	tmp := strings.Join(str, ",")
	tmp = strings.ReplaceAll(tmp, ",", "\",\"")
	return tmp
}

func readAppsFile(filePath string) ([]appgroupapp, error) {
	a := []appgroupapp{}

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

func getAPIProductsAsArray(prods []apiProduct) []string {
	var products []string
	for _, p := range prods {
		products = append(products, p.ApiProduct)
	}
	return products
}

func readAppsArrayFile(filePath string) ([][]appgroupapp, error) {
	a := [][]appgroupapp{}

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
