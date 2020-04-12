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

package apis

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	apiproxy "github.com/srinandan/apigeecli/cmd/apis/apiproxydef"
	proxies "github.com/srinandan/apigeecli/cmd/apis/proxies"
	target "github.com/srinandan/apigeecli/cmd/apis/targetendpoint"
)

func LoadDocument(filePath string) (doc *openapi3.Swagger, err error) {
	doc, err = openapi3.NewSwaggerLoader().LoadSwaggerFromFile(filePath)
	return doc, err
}

func GenerateAPIProxyDefFromOAS(name string, filePath string) (err error) {
	doc, err := LoadDocument(filePath)
	if err != nil {
		return err
	}

	apiproxy.SetDisplayName(name)
	if doc.Info != nil {
		if doc.Info.Description != "" {
			apiproxy.SetDescription(doc.Info.Description)
		}
	}

	apiproxy.SetCreatedAt()
	apiproxy.SetLastModifiedAt()
	apiproxy.SetConfigurationVersion()
	apiproxy.AddTargetEndpoint("default")
	apiproxy.AddProxyEndpoint("default")

	u, err := GetEndpoint(doc)
	if err != nil {
		return err
	}

	apiproxy.SetBasePath(u.Path)

	target.NewTargetEndpoint(u.Scheme + "://" + u.Hostname())

	proxies.NewProxyEndpoint(u.Path)

	GenerateFlows(doc.Paths)

	return nil
}

func GetEndpoint(doc *openapi3.Swagger) (u *url.URL, err error) {
	if doc.Servers == nil {
		return nil, fmt.Errorf("at least one server must be present")
	}

	return url.Parse(doc.Servers[0].URL)
}

func GetHTTPMethod(pathItem *openapi3.PathItem, keyPath string) map[string]string {

	pathMap := make(map[string]string)
	alternateOperationId := strings.ReplaceAll(keyPath, "\\", "_")

	if pathItem.Get != nil {
		if pathItem.Get.OperationID != "" {
			pathMap["get"] = pathItem.Get.OperationID
		} else {
			pathMap["get"] = "get_" + alternateOperationId
		}
	}

	if pathItem.Post != nil {
		if pathItem.Post.OperationID != "" {
			pathMap["post"] = pathItem.Post.OperationID
		} else {
			pathMap["post"] = "post_" + alternateOperationId
		}
	}

	if pathItem.Put != nil {
		if pathItem.Put.OperationID != "" {
			pathMap["put"] = pathItem.Put.OperationID
		} else {
			pathMap["put"] = "put_" + alternateOperationId
		}
	}

	if pathItem.Patch != nil {
		if pathItem.Patch.OperationID != "" {
			pathMap["patch"] = pathItem.Patch.OperationID
		} else {
			pathMap["patch"] = "patch_" + alternateOperationId
		}
	}

	if pathItem.Delete != nil {
		if pathItem.Delete.OperationID != "" {
			pathMap["delete"] = pathItem.Delete.OperationID
		} else {
			pathMap["delete"] = "delete_" + alternateOperationId
		}
	}

	if pathItem.Options != nil {
		if pathItem.Options.OperationID != "" {
			pathMap["options"] = pathItem.Options.OperationID
		} else {
			pathMap["options"] = "options_" + alternateOperationId
		}
	}

	if pathItem.Trace != nil {
		if pathItem.Trace.OperationID != "" {
			pathMap["trace"] = pathItem.Trace.OperationID
		} else {
			pathMap["trace"] = "trace_" + alternateOperationId
		}
	}

	if pathItem.Head != nil {
		if pathItem.Head.OperationID != "" {
			pathMap["head"] = pathItem.Head.OperationID
		} else {
			pathMap["head"] = "head_" + alternateOperationId
		}
	}

	return pathMap
}

func GenerateFlows(paths openapi3.Paths) {
	for keyPath := range paths {
		pathMap := GetHTTPMethod(paths[keyPath], keyPath)
		for method, operationId := range pathMap {
			proxies.AddFlow(operationId, replacePathWithWildCard(keyPath), method)
		}
	}
}

func replacePathWithWildCard(keyPath string) string {
	re := regexp.MustCompile(`{(.*?)}`)
	if strings.ContainsAny(keyPath, "{") {
		return re.ReplaceAllLiteralString(keyPath, "*")
	} else {
		return keyPath
	}
}
