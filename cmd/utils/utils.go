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

package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"internal/clilog"
)

func ListKVMFiles(folder string) (orgKVMFileList map[string]string,
	envKVMFileList map[string]string, proxyKVMFileList map[string]string, err error,
) {
	orgKVMFileList = map[string]string{}
	envKVMFileList = map[string]string{}
	proxyKVMFileList = map[string]string{}

	renv := regexp.MustCompile(`env_(\S*)_kvmfile_[0-9]+\.json`)
	rorg := regexp.MustCompile(`org_(\S*)_kvmfile_[0-9]+\.json`)
	rproxy := regexp.MustCompile(`proxy_(\S*)_kvmfile_[0-9]+\.json`)

	err = filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			kvmFile := filepath.Base(path)
			switch {
			case renv.MatchString(kvmFile):
				envKVMFileSplit := strings.Split(kvmFile, "_")
				if len(envKVMFileSplit) > 2 {
					clilog.Info.Printf("Map name %s, path %s\n", envKVMFileSplit[2], kvmFile)
					envKVMFileList[envKVMFileSplit[2]] = path
				}
			case rproxy.MatchString(kvmFile):
				proxyKVMFileSplit := strings.Split(kvmFile, "_")
				if len(proxyKVMFileSplit) > 2 {
					clilog.Info.Printf("Map name %s, path %s\n", proxyKVMFileSplit[2], kvmFile)
					proxyKVMFileList[proxyKVMFileSplit[2]] = path
				}
			case rorg.MatchString(kvmFile):
				orgKVMFileSplit := strings.Split(kvmFile, "_")
				if len(orgKVMFileSplit) > 1 {
					clilog.Info.Printf("Map name %s, path %s\n", orgKVMFileSplit[1], kvmFile)
					orgKVMFileList[orgKVMFileSplit[1]] = path
				}
			}
		}
		return nil
	})
	return orgKVMFileList, envKVMFileList, proxyKVMFileList, err
}

func ReadEntityFile(filePath string) ([]string, error) {
	entities := []string{}

	jsonFile, err := os.Open(filePath)
	if err != nil {
		return entities, err
	}

	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return entities, err
	}

	if err = json.Unmarshal(byteValue, &entities); err != nil {
		return entities, err
	}

	return entities, nil
}

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadFile(filePath string) (byteValue []byte, err error) {
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
