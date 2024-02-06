// Copyright 2024 Google LLC
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

package bundlegen

import (
	"testing"
)

func TestGenerateAPIProxyDefFromGQL(t *testing.T) {
	skipPolicy := false
	name := "test"
	gqlDocName := "schema"
	apiKeyLocation := "request.header.x-api-key"
	addCORS := true
	basePath := ""
	targetURL := "http://api.example.com"
	targetURLRef := ""
	if err := GenerateAPIProxyDefFromGQL(name,
		gqlDocName,
		basePath,
		apiKeyLocation,
		skipPolicy,
		addCORS,
		targetURLRef,
		targetURL); err != nil {
		t.Fatalf("%v", err)
	}
}
