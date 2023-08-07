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

package targetservers

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
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"

	"internal/clilog"
)

type targetserver struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Host        string   `json:"host,omitempty"`
	Port        int      `json:"port,omitempty"`
	IsEnabled   *bool    `json:"isEnabled,omitempty"`
	Protocol    string   `json:"protocol,omitempty"`
	SslInfo     *sslInfo `json:"sSLInfo,omitempty"`
}

type sslInfo struct {
	Enabled                *bool       `json:"enabled,omitempty"`
	ClientAuthEnabled      *bool       `json:"clientAuthEnabled,omitempty"`
	Keystore               string      `json:"keyStore,omitempty"`
	Keyalias               string      `json:"keyAlias,omitempty"`
	Truststore             string      `json:"trustStore,omitempty"`
	IgnoreValidationErrors *bool       `json:"ignoreValidationErrors,omitempty"`
	Protocols              []string    `json:"protocols,omitempty"`
	Ciphers                []string    `json:"ciphers,omitempty"`
	CommonName             *commonName `json:"commonName,omitempty"`
}

type commonName struct {
	Value         string `json:"value,omitempty"`
	WildcardMatch bool   `json:"wildcardMatch,omitempty"`
}

// Create
func Create(name string, description string, host string, port int, enable bool, protocol string, keyStore string, keyAlias string, trustStore string, tlsenabled string, clientAuthEnabled string, ignoreValidationErrors string) (respBody []byte, err error) {
	e := new(bool)
	*e = enable

	targetsvr := targetserver{
		Name:      name,
		IsEnabled: e,
	}

	return createOrUpdate("create", targetsvr, name, description, host, port, protocol, keyStore, keyAlias, trustStore, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
}

// Update
func Update(name string, description string, host string, port int, enable bool, protocol string, keyStore string, keyAlias string, trustStore string, tlsenabled string, clientAuthEnabled string, ignoreValidationErrors string) (respBody []byte, err error) {
	apiclient.ClientPrintHttpResponse.Set(false)
	targetRespBody, err := Get(name)
	if err != nil {
		return nil, err
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	targetsvr := targetserver{}
	if err = json.Unmarshal(targetRespBody, &targetsvr); err != nil {
		return nil, err
	}

	targetsvr.IsEnabled = &enable

	return createOrUpdate("update", targetsvr, name, description, host, port, protocol, keyStore, keyAlias, trustStore, tlsenabled, clientAuthEnabled, ignoreValidationErrors)
}

func createOrUpdate(action string, targetsvr targetserver, name string, description string, host string, port int, protocol string, keyStore string, keyAlias string, trustStore string, tlsenabled string, clientAuthEnabled string, ignoreValidationErrors string) (respBody []byte, err error) {
	if description != "" {
		targetsvr.Description = description
	}
	if host != "" {
		targetsvr.Host = host
	}

	if port != -1 {
		targetsvr.Port = port
	}
	if protocol != "" {
		targetsvr.Protocol = protocol
	}

	if keyStore != "" || keyAlias != "" || trustStore != "" || tlsenabled != "" ||
		clientAuthEnabled != "" || ignoreValidationErrors != "" {
		if targetsvr.SslInfo == nil {
			targetsvr.SslInfo = &sslInfo{}
		}
		targetsvr.SslInfo.Keystore = keyStore
		targetsvr.SslInfo.Keyalias = keyAlias
		targetsvr.SslInfo.Truststore = trustStore
		if tlsenabled != "" {
			targetsvr.SslInfo.Enabled = getBool(tlsenabled)
		}
		if clientAuthEnabled != "" {
			targetsvr.SslInfo.ClientAuthEnabled = getBool(clientAuthEnabled)
		}
		if ignoreValidationErrors != "" {
			targetsvr.SslInfo.IgnoreValidationErrors = getBool(ignoreValidationErrors)
		}
	}

	reqBody, err := json.Marshal(targetsvr)
	if err != nil {
		return nil, err
	}

	u, _ := url.Parse(apiclient.BaseURL)
	if action == "create" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
		respBody, err = apiclient.HttpClient(u.String(), string(reqBody))
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
		respBody, err = apiclient.HttpClient(u.String(), string(reqBody), "PUT")
	}

	return respBody, err
}

// Get
func Get(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// Delete
func Delete(name string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers", name)
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// List
func List() (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
	respBody, err = apiclient.HttpClient(u.String())
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

	var targetservers []string
	err = json.Unmarshal(respBody, &targetservers)
	if err != nil {
		return nil, err
	}

	clilog.Debug.Printf("Found %d targetservers in the org\n", len(targetservers))
	clilog.Debug.Printf("Exporting targetservers with %d parallel connections\n", conn)

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
		go exportServers(&fanOutWg, jobChan, resultChan, errChan)
	}

	for _, ts := range targetservers {
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

func exportServers(wg *sync.WaitGroup, jobs <-chan string, results chan<- []byte, errs chan<- error) {
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
	targetservers, err := readTargetServersFile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	clilog.Debug.Printf("Found %d target servers in the file\n", len(targetservers))
	clilog.Debug.Printf("Create target servers with %d connections\n", conn)

	u, _ := url.Parse(apiclient.BaseURL)
	u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
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

	var svrs []string
	if err = json.Unmarshal(b, &svrs); err != nil {
		return err
	}
	knownServers := map[string]bool{}
	for _, name := range svrs {
		knownServers[name] = true
	}

	jobChan := make(chan targetserver)
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
		go importServers(knownServers, &fanOutWg, jobChan, errChan)
	}

	for _, ts := range targetservers {
		jobChan <- ts
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

func importServers(knownServers map[string]bool, wg *sync.WaitGroup, jobs <-chan targetserver, errs chan<- error) {
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
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "targetservers")
		method := http.MethodPost
		if knownServers[job.Name] {
			// If the targetserver already exists we use an 'update' instead of a 'create' operation.
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
			errs <- fmt.Errorf("could not import targetserver, apigee responded with HTTP %d: %s", resp.StatusCode, resp.Status)
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
		clilog.Debug.Printf("Completed targetserver: %s", job.Name)
	}
}

func readTargetServersFile(filePath string) ([]targetserver, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var targetservers []targetserver
	if err = json.Unmarshal(content, &targetservers); err != nil {
		return nil, err
	}
	return targetservers, nil
}

func getBool(s string) *bool {
	tmp, _ := strconv.ParseBool(s)
	t := new(bool)
	t = &tmp
	return t
}
