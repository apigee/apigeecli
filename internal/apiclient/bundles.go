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

package apiclient

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"internal/clilog"
)

// entityPayloadList stores list of entities
var entityPayloadList [][]byte // types.EntityPayloadList

// ReadArchive confirms f the file format is zip and reads the contents are a byte[]
func ReadArchive(filename string) ([]byte, error) {
	if !strings.HasSuffix(filename, ".zip") {
		clilog.Error.Println("proxy bundle must be a zip file")
		return nil, errors.New("source must be a zipfile")
	}

	file, err := os.Open(filename)
	if err != nil {
		clilog.Error.Println("cannot open/read archive: ", err)
		return nil, err
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		clilog.Error.Println("error accessing file: ", err)
		return nil, err
	}

	_, err = zip.NewReader(file, fi.Size())

	if err != nil {
		clilog.Error.Println("invalid archive format: ", err)
		return nil, err
	}

	archiveFile, err := os.ReadFile(filename)
	if err != nil {
		clilog.Error.Printf("Error reading archive: %s\n", err)
		return nil, err
	}

	return archiveFile, nil
}

// ReadBundle confirms if the file format is a zip file
func ReadBundle(filename string) error {
	if !strings.HasSuffix(filename, ".zip") {
		clilog.Error.Println("proxy bundle must be a zip file")
		return errors.New("source must be a zipfile")
	}

	file, err := os.Open(filename)
	if err != nil {
		clilog.Error.Println("cannot open/read API Proxy Bundle: ", err)
		return err
	}

	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		clilog.Error.Println("error accessing file: ", err)
		return err
	}

	_, err = zip.NewReader(file, fi.Size())

	if err != nil {
		clilog.Error.Println("invalid API Proxy Bundle: ", err)
		return err
	}

	return nil
}

// WriteByteArrayToFile accepts []bytes and writes to a file
func WriteByteArrayToFile(exportFile string, fileAppend bool, payload []byte) error {
	fileFlags := os.O_CREATE | os.O_WRONLY

	if fileAppend {
		fileFlags |= os.O_APPEND
	} else {
		fileFlags |= os.O_TRUNC
	}

	f, err := os.OpenFile(exportFile, fileFlags, 0o644)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(payload)
	if err != nil {
		clilog.Error.Println("error writing to file: ", err)
		return err
	}
	return nil
}

// WriteArrayByteArrayToFile accepts [][]bytes and writes to a file
func WriteArrayByteArrayToFile(exportFile string, fileAppend bool, payload [][]byte) error {
	fileFlags := os.O_CREATE | os.O_WRONLY

	if fileAppend {
		fileFlags |= os.O_APPEND
	}

	f, err := os.OpenFile(exportFile, fileFlags, 0o644)
	if err != nil {
		return err
	}

	defer f.Close()

	// begin json array
	_, err = f.Write([]byte("["))
	if err != nil {
		clilog.Error.Println("error writing to file ", err)
		return err
	}

	payloadFromArray := bytes.Join(payload, []byte(","))
	// add json array terminate
	payloadFromArray = append(payloadFromArray, byte(']'))

	_, err = f.Write(payloadFromArray)

	if err != nil {
		clilog.Error.Println("error writing to file: ", err)
		return err
	}

	return nil
}

// GetAsyncEntity stores results for each entity in a list
func GetAsyncEntity(entityURL string, wg *sync.WaitGroup, mu *sync.Mutex) {
	// this is a two step process - 1) get entity details 2) store in byte[][]
	defer wg.Done()

	var respBody []byte

	// don't print to sysout
	ClientPrintHttpResponse.Set(false)
	defer ClientPrintHttpResponse.Set(GetCmdPrintHttpResponseSetting())
	respBody, err := HttpClient(entityURL)
	if err != nil {
		clilog.Error.Fatalf("error with entity: %s", entityURL)
		clilog.Error.Println(err)

		return
	}

	mu.Lock()
	entityPayloadList = append(entityPayloadList, respBody)
	mu.Unlock()
	clilog.Debug.Printf("Completed entity: %s", entityURL)
}

func GetEntityPayloadList() [][]byte {
	return entityPayloadList
}

func ClearEntityPayloadList() {
	entityPayloadList = entityPayloadList[:0]
}

// FetchAsyncBundle can download a shared flow or a proxy bundle
func FetchAsyncBundle(entityType string, folder string, name string, revision string, allRevisions bool, wg *sync.WaitGroup) {
	// this method is meant to be called asynchronously
	defer wg.Done()

	_ = FetchBundle(entityType, folder, name, revision, allRevisions)
}

// FetchBundle can download a shared flow or proxy bundle
func FetchBundle(entityType string, folder string, name string, revision string, allRevisions bool) error {
	var proxyName string

	u, _ := url.Parse(BaseURL)
	q := u.Query()
	q.Set("format", "bundle")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, GetApigeeOrg(), entityType, name, "revisions", revision)

	if allRevisions {
		proxyName = name + "_" + revision
	} else {
		proxyName = name
	}

	err := DownloadResource(u.String(), proxyName, ".zip", true)
	if err != nil {
		clilog.Error.Fatalf("error with entity: %s", name)
		clilog.Error.Println(err)
		return err
	}

	if len(folder) > 0 {
		_ = os.Rename(proxyName+".zip", path.Join(folder, proxyName+".zip"))
	}

	return nil
}

// ImportBundleAsync imports a sharedflow or api proxy bundle meantot be called asynchronously
func ImportBundleAsync(entityType string, name string, bundlePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	_ = ImportBundle(entityType, name, bundlePath)
}

// ImportBundle imports a sharedflow or api proxy bundle
func ImportBundle(entityType string, name string, bundlePath string) error {
	err := ReadBundle(bundlePath)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	// when importing from a folder, proxy name = file name
	if name == "" {
		_, fileName := filepath.Split(bundlePath)
		names := strings.Split(fileName, ".")
		name = names[0]
	}

	u, _ := url.Parse(BaseURL)
	u.Path = path.Join(u.Path, GetApigeeOrg(), entityType)

	q := u.Query()
	q.Set("name", name)
	q.Set("action", "import")
	u.RawQuery = q.Encode()

	formParams := map[string]string{
		"proxy": bundlePath,
	}

	_, err = PostHttpOctet(false, u.String(), formParams)
	if err != nil {
		clilog.Error.Println(err)
		return err
	}

	clilog.Debug.Printf("Completed entity: %s", u.String())
	return nil
}

func FolderExists(folder string) (err error) {
	if folder == "" {
		return nil
	}
	_, err = os.Stat(folder)
	if err != nil {
		return fmt.Errorf("folder not found or write permission denied")
	}
	return nil
}
