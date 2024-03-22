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

	apiclient "internal/apiclient"
	"internal/cmd"
	apicategories "internal/cmd/apicategories"
	products "internal/cmd/products"

	"github.com/spf13/cobra/doc"
)

var samples = `# apigeecli Samples
Set up apigeecli with preferences: apigeecli prefs set -o $org -t $(gcloud auth print-access-token)

| Operations | Command |
|---|---|
| apicategories | ` + apicategories.GetExample(0) + `|
| apis | apigeecli apis import -f samples/apis -o $org -t $token |
| appgroups | apigeecli appgroups import -f samples/appgroups.json -o $org -t $token |
| datacollectors | apigeecli datacollectors import -f samples/datacollectors.json -o $org -t $token |
| developers | apigeecli developers import -f samples/developers.json -o $org -t $token |
| kvms | apigeecli kvms import -f samples/kvms -o $org -t $token (Rename the files under samples/kvms to match your Apigee setup)  |
| products | ` + products.GetExample(2) + ` |
| products | ` + products.GetExample(3) + ` |
| products | ` + products.GetExample(4) + ` |
| products | ` + products.GetExample(0) + ` |
| products | ` + products.GetExample(1) + ` |
| sharedflows | apigeecli sharedflows import -f samples/sharedflows -o $org -t $token |
| targetservers | apigeecli targetservers import -f samples/targetservers.json -o $org -t $token -e $env |
| keystores | apigeecli keystores import -f samples/keystores.json -o $org -t $token -e $env |
| references | apigeecli references import -f samples/references.json -o $org -t $token -e $env |
| apps | apigeecli apps import -f samples/references.json -d samples/developers.json -o $org -t $token (better used when the developers and apps are exported using the export) |
| apidocs | apigeecli apidocs import -f samples/apidocs -o $org -t $token -s $siteId (Rename the files under samples/apidocs to match your Apigee setup for siteId) |`

func main() {
	var err error
	var docFiles []string

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
