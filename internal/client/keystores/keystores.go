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

package keystores

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"
)

// Create
func Create(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	payload := "{\"name\":\"" + name + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Import
func Import(conn int, filePath string) (err error) {
	keystores, err := readKeystoresFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Debug.Printf("Found %d keystores in the file\n", len(keystores))
	clilog.Debug.Printf("Create keystores with %d connections\n", conn)

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
		go importKeystores(&fanOutWg, jobChan, errChan)
	}

	for _, ks := range keystores {
		jobChan <- ks
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

func importKeystores(wg *sync.WaitGroup, jobs <-chan string, errs chan<- error) {
	defer wg.Done()

	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keystores")
		u.RawQuery = fmt.Sprintf("name=%s", job)

		err := apiclient.GetHttpClient()
		if err != nil {
			errs <- err
			continue
		}

		req, err := http.NewRequest(http.MethodPost, u.String(), nil)
		if err != nil {
			errs <- err
			continue
		}
		req, err = apiclient.SetAuthHeader(req)
		if err != nil {
			errs <- err
			continue
		}

		resp, err := apiclient.ApigeeAPIClient.Do(req)
		if err != nil {
			errs <- err
			continue
		} else if resp.StatusCode/100 != 2 && resp.StatusCode != http.StatusConflict {
			// We ignore 409s as the only configurable parameter of a keystore is it's name. Hence if it already exists
			// then it is consistent with what is being imported.
			errs <- fmt.Errorf("apigee responded with HTTP %d: %s", resp.StatusCode, resp.Status)
			continue
		}
	}
}

func readKeystoresFile(filePath string) ([]string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var keystoresList []string
	if err = json.Unmarshal(b, &keystoresList); err != nil {
		return nil, err
	}
	return keystoresList, nil
}
