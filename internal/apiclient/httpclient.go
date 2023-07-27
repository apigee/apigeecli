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

package apiclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"internal/clilog"

	"golang.org/x/time/rate"
)

// RateLimitedHttpClient
type RateLimitedHTTPClient struct {
	client      *http.Client
	Ratelimiter *rate.Limiter
}

// RateLimitedHTTPClient for Apigee API
var ApigeeAPIClient *RateLimitedHTTPClient

// allow 1 every 100 milliseconds
var apigeeAPIRateLimit = rate.NewLimiter(rate.Every(100*time.Millisecond), 1)

// source: https://cloud.google.com/apigee/docs/api-platform/reference/limits#apigee-apis

// allow 1 every 1 second
var apigeeAnalyticsAPIRateLimit = rate.NewLimiter(rate.Every(time.Second), 1)

// source: https://cloud.google.com/apigee/docs/api-platform/reference/limits#analytics-apis

// disable rate limit
var noAPIRateLimit = rate.NewLimiter(rate.Inf, 1)

// PostHttpZip method is used to send resources, proxy bundles, shared flows etc.
func PostHttpZip(auth bool, method string, url string, headers map[string]string, zipfile string) (err error) {
	var req *http.Request

	payload, err := os.ReadFile(zipfile)
	if err != nil {
		return err
	}

	err = GetHttpClient()
	if err != nil {
		return err
	}

	if DryRun() {
		return nil
	}

	clilog.Debug.Println("Connecting to : ", url)
	req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return err
	}

	for headerName, headerValue := range headers {
		clilog.Debug.Printf("%s : %s\n", headerName, headerValue)
		req.Header.Set(headerName, headerValue)
	}

	if auth { // do not pass auth header when using with archives
		req, err = SetAuthHeader(req)
		if err != nil {
			return err
		}
	}

	resp, err := ApigeeAPIClient.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	return nil
}

// PostHttpOctet method is used to send resources, proxy bundles, shared flows etc.
func PostHttpOctet(update bool, url string, formParams map[string]string) (respBody []byte, err error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for formName, formParam := range formParams {
		file, err := os.Open(formParam)
		if err != nil {
			clilog.Error.Printf("failed to open the file %s with error: %v", formParam, err)
			return nil, err
		}
		// get filenanme without extension
		fileNameWithExt, _ := filepath.Abs(formParam)
		formValue := strings.TrimSuffix(fileNameWithExt, filepath.Ext(formParam))
		part, err := writer.CreateFormFile(formName, formValue)
		if err != nil {
			clilog.Error.Println("Error writing multi-part: ", err)
			return nil, err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			clilog.Error.Println("error copying multi-part: ", err)
			return nil, err
		}

		file.Close()
	}

	err = writer.Close()
	if err != nil {
		clilog.Error.Println("error closing multi-part: ", err)
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	var req *http.Request

	err = GetHttpClient()
	if err != nil {
		return nil, err
	}

	clilog.Debug.Println("Connecting to : ", url)
	if !update {
		req, err = http.NewRequest("POST", url, body)
	} else {
		req, err = http.NewRequest("PUT", url, body)
	}

	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	req, err = SetAuthHeader(req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := ApigeeAPIClient.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	return handleResponse(resp)
}

func DownloadFile(url string, auth bool) (resp *http.Response, err error) {
	err = GetHttpClient()
	if err != nil {
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	clilog.Debug.Println("Connecting to : ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	if auth {
		req, err = SetAuthHeader(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err = ApigeeAPIClient.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	} else if resp.StatusCode > 299 {
		clilog.Error.Printf("error in response, status %d: %s", resp.StatusCode, resp.Body)
		return nil, errors.New("error in response")
	}

	if resp == nil {
		clilog.Error.Println("error in response: Response was null")
		return nil, fmt.Errorf("error in response: Response was null")
	}
	return resp, err
}

// DownloadResource method is used to download resources, proxy bundles, sharedflows
func DownloadResource(url string, name string, resType string, auth bool) error {
	var filename string

	if resType == ".zip" {
		filename = name + ".zip"
	} else {
		filename = name
	}

	if DryRun() {
		return nil
	}

	out, err := os.Create(filename)
	if err != nil {
		clilog.Error.Println("error creating file: ", err)
		return err
	}
	defer out.Close()

	resp, err := DownloadFile(url, auth)
	if err != nil {
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		clilog.Error.Println("error writing response to file: ", err)
		return err
	}

	clilog.Info.Println("Resource " + filename + " completed")
	return nil
}

// HttpClient method is used to GET,POST,PUT or DELETE JSON data
func HttpClient(params ...string) (respBody []byte, err error) {
	// The first parameter instructs whether the output should be printed
	// The second parameter is url. If only one parameter is sent, assume GET
	// The third parameter is the payload. The two parameters are sent, assume POST
	// THe fourth parameter is the method. If three parameters are sent, assume method in param
	// The fifth parameter is content type
	var req *http.Request
	contentType := "application/json"

	err = GetHttpClient()
	if err != nil {
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	clilog.Debug.Println("Connecting to: ", params[0])

	switch paramLen := len(params); paramLen {
	case 1:
		req, err = http.NewRequest("GET", params[0], nil)
	case 2:
		clilog.Debug.Println("Payload: ", params[1])
		req, err = http.NewRequest("POST", params[0], bytes.NewBuffer([]byte(params[1])))
	case 3:
		if req, err = getRequest(params); err != nil {
			return nil, err
		}
	case 4:
		if req, err = getRequest(params); err != nil {
			return nil, err
		}
		contentType = params[3]
	default:
		return nil, errors.New("unsupported method")
	}

	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	req, err = SetAuthHeader(req)
	if err != nil {
		return nil, err
	}

	clilog.Debug.Println("Content-Type : ", contentType)
	req.Header.Set("Content-Type", contentType)

	resp, err := ApigeeAPIClient.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	return handleResponse(resp)
}

// PrettyPrint method prints formatted json
func PrettyPrint(contentType string, body []byte) error {
	if GetCmdPrintHttpResponseSetting() && ClientPrintHttpResponse.Get() {
		var prettyJSON bytes.Buffer
		// pretty print only json responses with body
		if strings.Contains(contentType, "json") && len(body) > 0 {
			err := json.Indent(&prettyJSON, body, "", "\t")
			if err != nil {
				clilog.Error.Println("error parsing response: ", err)
				return err
			}

			clilog.HttpResponse.Println(prettyJSON.String())
		}
	}
	return nil
}

// PrettifyJSON
func PrettifyJSON(body []byte) ([]byte, error) {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		clilog.Error.Println("error parsing response: ", err)
		return nil, err
	}
	return prettyJSON.Bytes(), nil
}

// Do the HTTP request
func (c *RateLimitedHTTPClient) Do(req *http.Request) (*http.Response, error) {
	ctx := context.Background()
	// Wait until the rate is below Apigee limits
	err := c.Ratelimiter.Wait(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetHttpClient returns new http client with a rate limiter
func GetHttpClient() (err error) {
	var apiRateLimit *rate.Limiter

	switch r := GetRate(); r {
	case ApigeeAPI:
		apiRateLimit = apigeeAPIRateLimit
	case ApigeeAnalyticsAPI:
		apiRateLimit = apigeeAnalyticsAPIRateLimit
	case None:
		apiRateLimit = noAPIRateLimit
	default:
		apiRateLimit = noAPIRateLimit
	}

	if GetProxyURL() != "" {
		if proxyUrl, err := url.Parse(GetProxyURL()); err != nil {
			ApigeeAPIClient = &RateLimitedHTTPClient{
				client: &http.Client{
					Transport: &http.Transport{
						Proxy: http.ProxyURL(proxyUrl),
					},
				},
				Ratelimiter: apiRateLimit,
			}
		} else {
			return err
		}
	} else {
		ApigeeAPIClient = &RateLimitedHTTPClient{
			client:      http.DefaultClient,
			Ratelimiter: apiRateLimit,
		}
	}
	return nil
}

func getRequest(params []string) (req *http.Request, err error) {
	if params[2] == "DELETE" {
		req, err = http.NewRequest("DELETE", params[0], nil)
	} else if params[2] == "PUT" {
		clilog.Debug.Println("Payload: ", params[1])
		req, err = http.NewRequest("PUT", params[0], bytes.NewBuffer([]byte(params[1])))
	} else if params[2] == "PATCH" {
		clilog.Debug.Println("Payload: ", params[1])
		req, err = http.NewRequest("PATCH", params[0], bytes.NewBuffer([]byte(params[1])))
	} else if params[2] == "POST" {
		clilog.Debug.Println("Payload: ", params[1])
		req, err = http.NewRequest("POST", params[0], bytes.NewBuffer([]byte(params[1])))
	} else {
		return nil, errors.New("unsupported method")
	}
	return req, err
}

func SetAuthHeader(req *http.Request) (*http.Request, error) {
	if GetApigeeToken() == "" {
		if err := SetAccessToken(); err != nil {
			return nil, err
		}
	}
	clilog.Debug.Println("Setting token : ", GetApigeeToken())
	req.Header.Add("Authorization", "Bearer "+GetApigeeToken())
	return req, nil
}

func handleResponse(resp *http.Response) (respBody []byte, err error) {
	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil {
		clilog.Error.Println("error in response: Response was null")
		return nil, nil
	}

	respBody, err = io.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return nil, err
	} else if resp.StatusCode > 399 {
		clilog.Debug.Printf("status code %d, error in response: %s\n", resp.StatusCode, string(respBody))
		clilog.HttpError.Println(string(respBody))
		return nil, errors.New(getErrorMessage(resp.StatusCode))
	}

	return respBody, PrettyPrint(resp.Header.Get("Content-Type"), respBody)
}

func getErrorMessage(statusCode int) string {
	switch statusCode {
	case 400:
		return "Bad Request - malformed request syntax"
	case 401:
		return "Unauthorized - the client must authenticate itself"
	case 403:
		return "Forbidden - the client does not have access rights"
	case 404:
		return "Not found - the server cannot find the requested resource"
	case 405:
		return "Method Not Allowed - the request method is not supported by the target resource"
	case 409:
		return "Conflict - request conflicts with the current state of the server"
	case 415:
		return "Unsupported media type - media format of the requested data is not supported by the server"
	case 429:
		return "Too Many Request - user has sent too many requests"
	case 500:
		return "Internal server error"
	case 501:
		return "Not Implemented - request method is not supported by the server"
	case 502:
		return "Bad Gateway"
	case 503:
		return "Service Unavaliable - the server is not ready to handle the request"
	default:
		return "unknown error"
	}
}
