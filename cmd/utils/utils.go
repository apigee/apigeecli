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

// Builds a list of KVM files.
func buildKVMFileList(refKVMFileList map[string]string, skipExistingKVMs bool,
	kvmFileSplit []string, kvmFileSplitIndex int, kvmFilePath string, path string, level string,
) (err error) {
	var existingKVMListBytes []byte
	var kvmListArg string

	if level == "Proxy" {
		kvmListArg = kvmFileSplit[1] // kvm.List function expects name of the Proxy in case the KVM import level is "Proxy".
	} else {
		kvmListArg = ""
	}

	if existingKVMListBytes, err = kvm.List(kvmListArg); err != nil {
		return err
	}

	if skipExistingKVMs {
		if contains, _ := IfListContainsKVM(existingKVMListBytes, kvmFileSplit[kvmFileSplitIndex]); contains {
			clilog.Info.Printf("Skipping Map name %s having path %s as it already exists at the Apigee %s level\n",
				kvmFileSplit[kvmFileSplitIndex], kvmFilePath, level)
		} else {
			clilog.Info.Printf("Map name %s, path %s\n", kvmFileSplit[kvmFileSplitIndex], kvmFilePath)
			refKVMFileList[kvmFileSplit[kvmFileSplitIndex]] = path
		}
	} else {
		clilog.Info.Printf("Map name %s, path %s\n", kvmFileSplit[kvmFileSplitIndex], kvmFilePath)
		refKVMFileList[kvmFileSplit[kvmFileSplitIndex]] = path
	}

	return nil
}

func ListKVMFiles(folder string, skipExistingKVMs bool) (orgKVMFileList map[string]string,
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
					apiclient.SetApigeeEnv(envKVMFileSplit[1])
					if err = buildKVMFileList(envKVMFileList, skipExistingKVMs,
						envKVMFileSplit, 2, kvmFile, path, "Environment"); err != nil {
						return err
					}
				}
			case rproxy.MatchString(kvmFile):
				proxyKVMFileSplit := strings.Split(kvmFile, "_")
				if len(proxyKVMFileSplit) > 2 {
					if err = buildKVMFileList(proxyKVMFileList, skipExistingKVMs,
						proxyKVMFileSplit, 2, kvmFile, path, "Proxy"); err != nil {
						return err
					}
				}
			case rorg.MatchString(kvmFile):
				orgKVMFileSplit := strings.Split(kvmFile, "_")
				if len(orgKVMFileSplit) > 1 {
					if err = buildKVMFileList(orgKVMFileList, skipExistingKVMs,
						orgKVMFileSplit, 1, kvmFile, path, "Organization"); err != nil {
						return err
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
