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
	"log"
	"os"
	"path/filepath"

	cmd "internal/cmd"

	"github.com/spf13/cobra/doc"
)

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
}
