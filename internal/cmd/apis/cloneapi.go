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

package apis

import (
	"archive/zip"
	"fmt"
	"internal/apiclient"
	"internal/client/apis"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	proxybundle "internal/bundlegen/proxybundle"

	"github.com/spf13/cobra"
)

var CloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone an API proxy from an Zip or folder",
	Long:  "Clone an API proxy from an Zip or folder",
	Args: func(cmd *cobra.Command, args []string) (err error) {
		if proxyZip != "" && proxyFolder != "" {
			return fmt.Errorf("proxy bundle (zip) and folder to an API proxy cannot be combined")
		}
		if proxyZip == "" && proxyFolder == "" {
			return fmt.Errorf("either proxy bundle (zip) or folder must be specified, not both")
		}
		if proxyFolder != "" {
			if _, err := os.Stat(proxyFolder); os.IsNotExist(err) {
				return err
			}
		}
		if !strings.HasPrefix(basePath, "/") {
			return fmt.Errorf("basePath must start with /")
		}
		apiclient.SetRegion(region)
		return apiclient.SetApigeeOrg(org)
	},
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		cmd.SilenceUsage = true

		if proxyZip != "" {
			// extract the zip to a tmp folder and assign to proxyFolder
			if proxyFolder, err = unzipBundle(); err != nil {
				return err
			}
		}

		if proxyFolder != "" {
			var tmpDir string

			if proxyZip != "" {
				tmpDir = proxyFolder
			} else {
				if tmpDir, err = copyDirectory(); err != nil {
					return err
				}
				defer os.RemoveAll(tmpDir)
			}

			if err = renameProxy(tmpDir); err != nil {
				return err
			}

			proxyBundlePath := path.Join(tmpDir, name+zipExt)

			if err = proxybundle.GenerateArchiveBundle(path.Join(tmpDir, "apiproxy"), proxyBundlePath, false); err != nil {
				return err
			}
			if _, err = apis.CreateProxy(name, proxyBundlePath); err != nil {
				return err
			}

			return err
		}

		return err
	},
}

func init() {
	CloneCmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name for the cloned proxy")

	CloneCmd.Flags().StringVarP(&proxyZip, "proxy-zip", "p",
		"", "Path to the Sharedflow bundle/zip file")
	CloneCmd.Flags().StringVarP(&basePath, "basepath", "b",
		"", "A new basePath for the cloned proxy")
	CloneCmd.Flags().StringVarP(&proxyFolder, "proxy-folder", "f",
		"", "Path to the Sharedflow Bundle; ex: ./test/apiproxy")

	_ = CloneCmd.MarkFlagRequired("name")
	_ = CloneCmd.MarkFlagRequired("basepath")
}

func copyDirectory() (tmpDir string, err error) {
	if strings.Contains(runtime.GOOS, "windows") {
		return "", fmt.Errorf("this operation is not supported on windows at the moment")
	}

	tmpDir, err = os.MkdirTemp("", "proxy")
	if err != nil {
		return tmpDir, err
	}

	copyCmd := exec.Command("cp", "-r", proxyFolder, tmpDir)
	if err = copyCmd.Run(); err != nil {
		return tmpDir, err
	}
	return tmpDir, nil
}

func renameProxy(tmpDir string) (err error) {
	// 1. rename the file in the apiproxy folder
	apiproxyFolder := path.Join(tmpDir, "apiproxy")
	re := regexp.MustCompile(`apiproxy\/\w+\.xml`)

	err = filepath.Walk(apiproxyFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			oldFileName := filepath.Base(path)

			if re.MatchString(path) { // this is the proxy xml
				// 1. rename the file based on the new proxy name
				newFilePath := strings.Replace(path, oldFileName, name+".xml", 1)
				if err = os.Rename(path, newFilePath); err != nil {
					return err
				}

				// 2. set the name
				if err = setParam(newFilePath, "proxy"); err != nil {
					return err
				}

				// 3. set the basePath
				if err = setParam(newFilePath, "basePath"); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return setBasePath(path.Join(tmpDir, "apiproxy", "proxies", "default.xml"))
}

func setParam(filePath string, paramType string) (err error) {
	var proxyFile *os.File

	proxyFile, err = os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}

	byteValue, err := io.ReadAll(proxyFile)
	if err != nil {
		return err
	}

	proxyFile.Close()

	stringValue := string(byteValue)
	replaceName := fmt.Sprintf("<APIProxy revision=\"1\" name=\"%s\">", name)
	replaceBasePath := fmt.Sprintf("<Basepaths>%s</Basepaths>", basePath)

	switch paramType {
	case "proxy":
		re := regexp.MustCompile(`\<APIProxy revision=\"\d+\" name=\"\S+\">`)
		stringValue = re.ReplaceAllString(stringValue, replaceName)
	case "basePath":
		re := regexp.MustCompile(`\<Basepaths>\/\w+\<\/Basepaths>`)
		stringValue = re.ReplaceAllString(stringValue, replaceBasePath)
	}

	proxyFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer proxyFile.Close()

	_, err = proxyFile.Write([]byte(stringValue))
	if err != nil {
		return err
	}

	return nil
}

func setBasePath(filePath string) (err error) {
	var proxyFile *os.File

	replaceBasePath := fmt.Sprintf("<BasePath>%s</BasePath>", basePath)

	proxyFile, err = os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return err
	}

	byteValue, err := io.ReadAll(proxyFile)
	if err != nil {
		return err
	}

	proxyFile.Close()

	stringValue := string(byteValue)

	re := regexp.MustCompile(`\<BasePath>\/\w+\<\/BasePath>`)
	stringValue = re.ReplaceAllString(stringValue, replaceBasePath)

	proxyFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}

	defer proxyFile.Close()

	_, err = proxyFile.Write([]byte(stringValue))
	if err != nil {
		return err
	}

	return nil
}

func unzipBundle() (tmpDir string, err error) {
	tmpDir, err = os.MkdirTemp("", "proxy")
	if err != nil {
		return tmpDir, err
	}

	bundle, err := zip.OpenReader(proxyZip)
	if err != nil {
		return tmpDir, err
	}

	for _, item := range bundle.File {
		if !strings.Contains(item.Name, "..") {

			bundlePath := filepath.Join(tmpDir, item.Name)

			if item.FileInfo().IsDir() {
				if err = os.MkdirAll(bundlePath, os.ModePerm); err != nil {
					return tmpDir, err
				}
				continue
			}

			if err = os.MkdirAll(filepath.Dir(bundlePath), os.ModePerm); err != nil {
				return tmpDir, err
			}

			bundleFile, err := os.OpenFile(bundlePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				return tmpDir, err
			}

			zipFile, err := item.Open()
			if err != nil {
				return tmpDir, err
			}

			if _, err = io.Copy(bundleFile, zipFile); err != nil {
				return tmpDir, err
			}

			bundleFile.Close()
			zipFile.Close()
		}
	}

	return tmpDir, nil
}
