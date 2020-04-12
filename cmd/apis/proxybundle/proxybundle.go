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

package proxybundle

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"

	apiproxy "github.com/srinandan/apigeecli/cmd/apis/apiproxydef"
	proxies "github.com/srinandan/apigeecli/cmd/apis/proxies"
	target "github.com/srinandan/apigeecli/cmd/apis/targetendpoint"
)

func GenerateAPIProxyBundle(name string) (err error) {
	const rootDir = "apiproxy"

	err = os.Mkdir(rootDir, os.ModePerm)
	if err != nil {
		return err
	}

	// write API Proxy file
	err = writeXMLData(rootDir+string(os.PathSeparator)+name+".xml", apiproxy.GetAPIProxy())
	if err != nil {
		return err
	}

	proxiesDirPath := rootDir + string(os.PathSeparator) + "proxies"
	targetDirPath := rootDir + string(os.PathSeparator) + "targets"

	err = os.Mkdir(proxiesDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = writeXMLData(proxiesDirPath+string(os.PathSeparator)+"default.xml", proxies.GetProxyEndpoint())
	if err != nil {
		return err
	}

	err = os.Mkdir(targetDirPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = writeXMLData(targetDirPath+string(os.PathSeparator)+"default.xml", target.GetTargetEndpoint())
	if err != nil {
		return err
	}

	err = archiveBundle(rootDir, name+".zip")
	if err != nil {
		return err
	}

	defer os.RemoveAll(rootDir) // clean up
	return nil
}

func writeXMLData(fileName string, data string) error {
	fileWriter, err := os.Create(fileName)
	if err != nil {
		return err
	}
	_, err = fileWriter.WriteString(data)
	if err != nil {
		return err
	}

	fileWriter.Close()
	return nil
}

func archiveBundle(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}
