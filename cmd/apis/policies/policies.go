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

package policies

import (
	"regexp"
)

var oasPolicyTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<OASValidation continueOnError="false" enabled="true" name="OpenAPI-Spec-Validation-1">
    <DisplayName>OpenAPI Spec Validation-1</DisplayName>
    <Properties/>
    <Source>request</Source>
    <OASResource>oas://{PolicyName}</OASResource>
</OASValidation>`

func AddOpenAPIValidatePolicy(name string) string {
	return replaceTemplateWithPolicy(name)
}

func replaceTemplateWithPolicy(name string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	return re.ReplaceAllLiteralString(oasPolicyTemplate, name)
}
