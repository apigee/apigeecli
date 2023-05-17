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

package references

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"
)

type ref struct {
	Name         string `json:"name,omitempty"`
	ResourceType string `json:"resourceType,omitempty"`
	Refers       string `json:"refers,omitempty"`
}

// Create references
func Create(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	reference = append(reference, "\"refers\":\""+refers+"\"")

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// Get a reference
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete a reference
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List references
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Update references
func Update(name string, description string, resourceType string, refers string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)

	reference := []string{}

	reference = append(reference, "\"name\":\""+name+"\"")

	if description != "" {
		reference = append(reference, "\"description\":\""+description+"\"")
	}

	if resourceType != "" {
		reference = append(reference, "\"resourceType\":\""+resourceType+"\"")
	}

	if refers != "" {
		reference = append(reference, "\"refers\":\""+refers+"\"")
	}

	payload := "{" + strings.Join(reference, ",") + "}"

	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references", name)
	respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	return respBody, err
}

// Export
func Export(conn int) (payload [][]byte, err error) {
	// don't print to sysout
	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	respBody, err := List()
	if err != nil {
		return nil, err
	}

	var references []string
	err = json.Unmarshal(respBody, &references)
	if err != nil {
		return nil, err
	}

	clilog.Debug.Printf("Found %d references in the org\n", len(references))
	clilog.Debug.Printf("Exporting references with %d connections\n", conn)

	jobChan := make(chan string)
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

	results := [][]byte{}
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
		go exportReferences(&fanOutWg, jobChan, resultChan, errChan)
	}

	for _, ts := range references {
		jobChan <- ts
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

func exportReferences(wg *sync.WaitGroup, jobs <-chan string, results chan<- []byte, errs chan<- error) {
	defer wg.Done()
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	apiclient.ClientPrintHttpResponse.Set(false)
	for {
		job, ok := <-jobs
		if !ok {
			return
		}
		respBody, err := Get(job)
		if err != nil {
			errs <- err
		} else {
			results <- respBody
		}
	}
}

// Import
func Import(conn int, filePath string) (err error) {
	references, err := readReferencesFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Debug.Printf("Found %d references in the file\n", len(references))
	clilog.Debug.Printf("Create references with %d connections\n", conn)

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
	err = apiclient.GetHttpClient()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	req, err = apiclient.SetAuthHeader(req)
	if err != nil {
		return err
	}

	resp, err := apiclient.ApigeeAPIClient.Do(req)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var refs []string
	if err = json.Unmarshal(b, &refs); err != nil {
		return err
	}
	knownRefs := map[string]bool{}
	for _, name := range refs {
		knownRefs[name] = true
	}

	jobChan := make(chan ref)
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
		go importReferences(knownRefs, &fanOutWg, jobChan, errChan)
	}

	for _, ref := range references {
		jobChan <- ref
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

func importReferences(knownRefs map[string]bool, wg *sync.WaitGroup, jobs <-chan ref, errs chan<- error) {
	defer wg.Done()

	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		b, err := json.Marshal(job)
		if err != nil {
			errs <- err
			continue
		}

		err = apiclient.GetHttpClient()
		if err != nil {
			errs <- err
			continue
		}

		u, _ := url.Parse(apiclient.BaseURL)
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "references")
		method := http.MethodPost
		if knownRefs[job.Name] {
			// If the reference already exists we use an 'update' instead of a 'create' operation.
			u.Path = path.Join(u.Path, job.Name)
			method = http.MethodPut
		}

		req, err := http.NewRequest(method, u.String(), bytes.NewReader(b))
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
			errs <- fmt.Errorf("failed to import reference, apigee returned HTTP %d: %s", resp.StatusCode, resp.Status)
			continue
		}

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			errs <- err
			continue
		}

		if len(b) > 0 && apiclient.GetPrintOutput() {
			out := bytes.NewBuffer([]byte{})
			if err = json.Indent(out, bytes.TrimSpace(b), "", "  "); err != nil {
				errs <- fmt.Errorf("apigee returned invalid json: %w", err)
			}
		}
		clilog.Debug.Printf("Completed reference: %s", job.Name)
	}
}

func readReferencesFile(filePath string) ([]ref, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var refList []ref
	if err = json.Unmarshal(b, &refList); err != nil {
		return nil, err
	}
	return refList, nil
}
