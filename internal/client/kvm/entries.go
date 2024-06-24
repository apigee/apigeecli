// Copyright 2021 Google LLC
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

package kvm

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"internal/apiclient"
	"internal/client/apis"

	"internal/clilog"
	"internal/cmd/utils"
)

type keyvalueentry struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type keyvalueentries struct {
	KeyValueEntries []keyvalueentry `json:"keyValueEntries,omitempty"`
	NextPageToken   string          `json:"nextPageToken,omitempty"`
}

// CreateEntry
func CreateEntry(proxyName string, mapName string, keyName string, value string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries")
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries")
	}
	payload := "{\"name\":\"" + keyName + "\",\"value\":\"" + value + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload)
	return respBody, err
}

// DeleteEntry
func DeleteEntry(proxyName string, mapName string, keyName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries", keyName)
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries", keyName)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries", keyName)
	}
	respBody, err = apiclient.HttpClient(u.String(), "", "DELETE")
	return respBody, err
}

// GetEntry
func GetEntry(proxyName string, mapName string, keyName string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries", keyName)
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries", keyName)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries", keyName)
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// UpdateEntry
func UpdateEntry(proxyName string, mapName string, keyName string, value string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())
	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries", keyName)
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries", keyName)
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries", keyName)
	}
	payload := "{\"name\":\"" + keyName + "\",\"value\":\"" + value + "\"}"
	respBody, err = apiclient.HttpClient(u.String(), payload, "PUT")
	return respBody, err
}

func upsertEntry(proxyName string, mapName string, keyName string, value string) (respBody []byte, err error) {
	update := false
	apiclient.ClientPrintHttpResponse.Set(false)
	_, err = GetEntry(proxyName, mapName, keyName)
	if err == nil {
		update = true
	}
	apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())
	if update {
		return UpdateEntry(proxyName, mapName, keyName, value)
	}
	return CreateEntry(proxyName, mapName, keyName, value)
}

// ListEntries
func ListEntries(proxyName string, mapName string, pageSize int, pageToken string) (respBody []byte, err error) {
	u, _ := url.Parse(apiclient.GetApigeeBaseURL())

	if pageToken != "" || pageSize != -1 {
		q := u.Query()
		if pageToken != "" {
			q.Set("page_token", pageToken)
		}
		if pageSize != -1 {
			q.Set("page_size", strconv.Itoa(pageSize))
		}
		u.RawQuery = q.Encode()
	}

	if apiclient.GetApigeeEnv() != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "environments", apiclient.GetApigeeEnv(), "keyvaluemaps", mapName, "entries")
	} else if proxyName != "" {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "apis", proxyName, "keyvaluemaps", mapName, "entries")
	} else {
		u.Path = path.Join(u.Path, apiclient.GetApigeeOrg(), "keyvaluemaps", mapName, "entries")
	}
	respBody, err = apiclient.HttpClient(u.String())
	return respBody, err
}

// ExportAllEntries
func ExportAllEntries() (err error) {
	type proxy struct {
		Name         string `json:"name,omitempty"`
		APIProxyType string `json:"apiProxyType,omitempty"`
	}

	type proxies struct {
		Proxies []proxy `json:"proxies,omitempty"`
	}

	p := proxies{}
	programmableProxies := []string{}

	apiList, err := apis.ListProxies(false)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(apiList, &p); err != nil {
		return err
	}

	for _, proxy := range p.Proxies {
		// search for only programmable proxies. standard proxies can't have KVMs
		if proxy.APIProxyType == "PROGRAMMABLE" {
			programmableProxies = append(programmableProxies, proxy.Name)
		}
	}

	for _, programmableProxy := range programmableProxies {
		kvmList := []string{}
		kvmListResp, err := List(programmableProxy)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(kvmListResp, &kvmList); err != nil {
			return err
		}
		if len(kvmList) > 0 {
			clilog.Info.Printf("Found %d scoped KVMs for proxy %s\n", len(kvmList), programmableProxy)
			for _, proxyKVM := range kvmList {
				clilog.Info.Printf("Exporting entries for KVM %s of proxy %s\n", proxyKVM, programmableProxy)
				proxyKVMEntries, err := ExportEntries(programmableProxy, proxyKVM)
				if err != nil {
					return err
				}
				fileName := strings.Join([]string{"proxy", programmableProxy, proxyKVM, "kvmfile"}, utils.DefaultFileSplitter)
				for i := range proxyKVMEntries {
					if err = apiclient.WriteByteArrayToFile(fileName+utils.DefaultFileSplitter+strconv.Itoa(i)+".json", false, proxyKVMEntries[i]); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// ExportEntries
func ExportEntries(proxyName string, mapName string) (payload [][]byte, err error) {
	var respBody []byte
	count := 1

	apiclient.ClientPrintHttpResponse.Set(false)
	defer apiclient.ClientPrintHttpResponse.Set(apiclient.GetCmdPrintHttpResponseSetting())

	if respBody, err = ListEntries(proxyName, mapName, -1, ""); err != nil {
		return nil, err
	}

	clilog.Debug.Printf("Exporting batch 1 of KVM entries for map %s\n", mapName)
	payload = append(payload, respBody)

	keyValueEntries := keyvalueentries{}
	if err = json.Unmarshal(respBody, &keyValueEntries); err != nil {
		return nil, err
	}

	for keyValueEntries.NextPageToken != "" {
		if respBody, err = ListEntries(proxyName, mapName, -1, keyValueEntries.NextPageToken); err != nil {
			return nil, err
		}

		if err = json.Unmarshal(respBody, &keyValueEntries); err != nil {
			return nil, err
		}
		count++
		clilog.Debug.Printf("Exporting batch %d of KVM entries for map %s\n", count, mapName)
		payload = append(payload, respBody)
	}

	return payload, nil
}

// ImportEntries
func ImportEntries(proxyName string, mapName string, conn int, filePath string) (err error) {

	kvmEntries, err := readKVMfile(filePath)
	if err != nil {
		clilog.Error.Println("Error reading file: ", err)
		return err
	}

	numEntities := len(kvmEntries.KeyValueEntries)
	clilog.Debug.Printf("Found %d entries in the file\n", numEntities)
	clilog.Debug.Printf("Create KVM entries with %d connections\n", conn)

	jobChan := make(chan keyvalueentry)
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
		go batchImport(&fanOutWg, proxyName, mapName, jobChan, errChan)
	}

	for _, entity := range kvmEntries.KeyValueEntries {
		jobChan <- entity
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

func batchImport(wg *sync.WaitGroup, proxyName string, mapName string, jobs <-chan keyvalueentry, errs chan<- error) {
	var err error
	defer wg.Done()
	for {
		job, ok := <-jobs
		if !ok {
			return
		}

		_, err = upsertEntry(proxyName, mapName, job.Name, job.Value)

		if err != nil {
			errs <- err
			continue
		}
	}
}

func readKVMfile(filePath string) (kvmEntries keyvalueentries, err error) {
	kvmEntries = keyvalueentries{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return kvmEntries, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return kvmEntries, err
	}

	err = json.Unmarshal(byteValue, &kvmEntries)

	if err != nil {
		return kvmEntries, err
	}

	return kvmEntries, nil
}
