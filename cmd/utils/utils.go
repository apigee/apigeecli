package utils

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"internal/apiclient"
	"internal/client/kvm"
	"internal/clilog"
)

// Checks if a specific KVM is present in the list of KVMs.
func IfListContainsKVM(refKVMListBytes []byte, kvm string) (containsKVM bool, err error) {
	var strRefKVMList []string

	if err = json.Unmarshal(refKVMListBytes, &strRefKVMList); err != nil {
		return false, err
	}

	for _, mapName := range strRefKVMList {
		if mapName == kvm {
			return true, nil
		}
	}

	return false, nil
}

func ListKVMFiles(folder string) (orgKVMFileList map[string]string,
	envKVMFileList map[string]string, proxyKVMFileList map[string]string, err error,
) {
	orgKVMFileList = map[string]string{}
	envKVMFileList = map[string]string{}
	proxyKVMFileList = map[string]string{}

	var existingKVMListBytes []byte

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
					apiclient.SetApigeeEnv(envKVMFileSplit[1])
					if existingKVMListBytes, err = kvm.List(""); err != nil {
						return err
					}
					if contains, _ := IfListContainsKVM(existingKVMListBytes, envKVMFileSplit[2]); contains {
						clilog.Info.Printf("Skipping Map name %s having path %s as it already exists at the Apigee Environment level\n",
							envKVMFileSplit[2], kvmFile)
					} else {
						clilog.Info.Printf("Map name %s, path %s\n", envKVMFileSplit[2], kvmFile)
						envKVMFileList[envKVMFileSplit[2]] = path
					}
				}
			case rproxy.MatchString(kvmFile):
				proxyKVMFileSplit := strings.Split(kvmFile, "_")
				if len(proxyKVMFileSplit) > 2 {
					if existingKVMListBytes, err = kvm.List(proxyKVMFileSplit[1]); err != nil {
						return err
					}
					if contains, _ := IfListContainsKVM(existingKVMListBytes, proxyKVMFileSplit[2]); contains {
						clilog.Info.Printf("Skipping Map name %s having path %s as it already exists at the Apigee API Proxy level\n",
							proxyKVMFileSplit[2], kvmFile)
					} else {
						clilog.Info.Printf("Map name %s, path %s\n", proxyKVMFileSplit[2], kvmFile)
						proxyKVMFileList[proxyKVMFileSplit[2]] = path
					}
				}
			case rorg.MatchString(kvmFile):
				orgKVMFileSplit := strings.Split(kvmFile, "_")
				if len(orgKVMFileSplit) > 1 {
					if existingKVMListBytes, err = kvm.List(""); err != nil {
						return err
					}
					if contains, _ := IfListContainsKVM(existingKVMListBytes, orgKVMFileSplit[1]); contains {
						clilog.Info.Printf("Skipping Map name %s having path %s as it already exists at the Apigee Organization level\n",
							orgKVMFileSplit[1], kvmFile)
					} else {
						clilog.Info.Printf("Map name %s, path %s\n", orgKVMFileSplit[1], kvmFile)
						orgKVMFileList[orgKVMFileSplit[1]] = path
					}
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
