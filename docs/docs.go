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

package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	apiclient "internal/apiclient"
	"internal/cmd"

	"github.com/spf13/cobra/doc"
)

const ENABLED = "true"

// var samples = `# apigeecli command Samples

// The following table contains some examples of apigeecli. Set up apigeecli with preferences: apigeecli prefs set -o $org

// | Operations | Command |
// |---|---|
// | apicategories | ` + getSingleLine(apicategories.GetExample(0)) + `|
// | apis | ` + getSingleLine(apis.GetExample(0)) + ` |
// | apis | ` + getSingleLine(apis.GetExample(1)) + ` |
// | apis | ` + getSingleLine(apis.GetExample(2)) + ` |
// | appgroups | ` + getSingleLine(appgroups.GetExample(0)) + ` |
// | datacollectors | ` + getSingleLine(datacollectors.GetExample(0)) + `  |
// | developers | ` + getSingleLine(developers.GetExample(0)) + `  |
// | kvms | ` + getSingleLine(kvm.GetExample(0)) + `  |
// | products | ` + getSingleLine(products.GetExample(2)) + ` |
// | products | ` + getSingleLine(products.GetExample(3)) + ` |
// | products | ` + getSingleLine(products.GetExample(4)) + ` |
// | products | ` + getSingleLine(products.GetExample(0)) + ` |
// | products | ` + getSingleLine(products.GetExample(1)) + ` |
// | sharedflows | ` + getSingleLine(sharedflows.GetExample(0)) + ` |
// | targetservers | ` + getSingleLine(targetservers.GetExample(0)) + ` |
// | keystores | `+getSingleLine(keystrores.GetExample(0)+` |
// | references | `+getSingleLine(references.GetExample(0))+` |
// | apps | `+getSingleLine(apps.GetExample(0))+` |
// | apidocs | `+getSingleLine(apidocs.GetExample(0))+`  |

// NOTE: This file is auto-generated during a release. Do not modify.`

var samples = `# apigeecli command Samples`

func main() {
	var err error
	var docFiles []string

	if os.Getenv("APIGEECLI_SKIP_DOCS") != ENABLED {

		if docFiles, err = filepath.Glob("./docs/apigeecli*.md"); err != nil {
			log.Fatal(err)
		}

		for _, docFile := range docFiles {
			if err = os.Remove(docFile); err != nil {
				log.Fatal(err)
			}
		}

		if err = doc.GenMarkdownTree(cmd.RootCmd, "./docs"); err != nil {
			log.Fatal(err)
		}
	}

	_ = apiclient.WriteByteArrayToFile("./samples/README.md", false, []byte(samples))
}

func WriteFile() (byteValue []byte, err error) {
	userFile, err := os.Open("./samples/README.md")
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

func getSingleLine(s string) string {
	return "`" + strings.ReplaceAll(strings.ReplaceAll(s, "\\", ""), "\n", "") + "`"
}
