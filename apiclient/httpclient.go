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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/srinandan/apigeecli/clilog"
)

//PostHttpZip method is used to send resources, proxy bundles, shared flows etc.
func PostHttpZip(print bool, auth bool, method string, url string, headers map[string]string, zipfile string) (err error) {

	var req *http.Request

	payload, err := ioutil.ReadFile(zipfile)
	if err != nil {
		return err
	}

	client, err := getHttpClient()
	if err != nil {
		return err
	}

	if DryRun() {
		return nil
	}

	clilog.Info.Println("Connecting to : ", url)
	req, err = http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return err
	}

	for headerName, headerValue := range headers {
		clilog.Info.Printf("%s : %s\n", headerName, headerValue)
		req.Header.Set(headerName, headerValue)
	}

	if auth { //do not pass auth header when using with archives
		req, err = setAuthHeader(req)
		if err != nil {
			return err
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	return nil
}

//PostHttpOctet method is used to send resources, proxy bundles, shared flows etc.
func PostHttpOctet(print bool, update bool, url string, proxyName string) (respBody []byte, err error) {
	file, err := os.Open(proxyName)
	if err != nil {
		clilog.Error.Printf("failed to open the file %s with error: %v", proxyName, err)
		return nil, err
	}
	defer file.Close()

	if DryRun() {
		return nil, nil
	}

	var req *http.Request

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("proxy", proxyName)
	if err != nil {
		clilog.Error.Println("Error writing multi-part: ", err)
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		clilog.Error.Println("error copying multi-part: ", err)
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		clilog.Error.Println("error closing multi-part: ", err)
		return nil, err
	}

	client, err := getHttpClient()
	if err != nil {
		return nil, err
	}

	clilog.Info.Println("Connecting to : ", url)
	if !update {
		req, err = http.NewRequest("POST", url, body)
	} else {
		req, err = http.NewRequest("PUT", url, body)
	}

	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	req, err = setAuthHeader(req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	return handleResponse(print, resp)
}

func DownloadFile(url string, auth bool) (resp *http.Response, err error) {
	client, err := getHttpClient()
	if err != nil {
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	clilog.Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	if auth {
		req, err = setAuthHeader(req)
		if err != nil {
			return nil, err
		}
	}

	resp, err = client.Do(req)

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

//DownloadResource method is used to download resources, proxy bundles, sharedflows
func DownloadResource(url string, name string, resType string) error {
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

	resp, err := DownloadFile(url, true)
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

	fmt.Println("Resource " + filename + " completed")
	return nil
}

//HttpClient method is used to GET,POST,PUT or DELETE JSON data
func HttpClient(print bool, params ...string) (respBody []byte, err error) {
	// The first parameter instructs whether the output should be printed
	// The second parameter is url. If only one parameter is sent, assume GET
	// The third parameter is the payload. The two parameters are sent, assume POST
	// THe fourth parameter is the method. If three parameters are sent, assume method in param
	//The fifth parameter is content type
	var req *http.Request
	contentType := "application/json"

	client, err := getHttpClient()
	if err != nil {
		return nil, err
	}

	if DryRun() {
		return nil, nil
	}

	clilog.Info.Println("Connecting to: ", params[0])

	switch paramLen := len(params); paramLen {
	case 1:
		req, err = http.NewRequest("GET", params[0], nil)
	case 2:
		clilog.Info.Println("Payload: ", params[1])
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

	req, err = setAuthHeader(req)
	if err != nil {
		return nil, err
	}

	clilog.Info.Println("Content-Type : ", contentType)
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	return handleResponse(print, resp)
}

//PrettyPrint method prints formatted json
func PrettyPrint(body []byte) error {
	var prettyJSON bytes.Buffer
	err := json.Indent(&prettyJSON, body, "", "\t")
	if err != nil {
		clilog.Error.Println("error parsing response: ", err)
		return err
	}
	fmt.Println(prettyJSON.String())
	return nil
}

func getRequest(params []string) (req *http.Request, err error) {
	if params[2] == "DELETE" {
		req, err = http.NewRequest("DELETE", params[0], nil)
	} else if params[2] == "PUT" {
		clilog.Info.Println("Payload: ", params[1])
		req, err = http.NewRequest("PUT", params[0], bytes.NewBuffer([]byte(params[1])))
	} else if params[2] == "PATCH" {
		clilog.Info.Println("Payload: ", params[1])
		req, err = http.NewRequest("PATCH", params[0], bytes.NewBuffer([]byte(params[1])))
	} else if params[2] == "POST" {
		clilog.Info.Println("Payload: ", params[1])
		req, err = http.NewRequest("POST", params[0], bytes.NewBuffer([]byte(params[1])))
	} else {
		return nil, errors.New("unsupported method")
	}
	return req, err
}

func setAuthHeader(req *http.Request) (*http.Request, error) {
	if GetApigeeToken() == "" {
		if err := SetAccessToken(); err != nil {
			return nil, err
		}
	}
	clilog.Info.Println("Setting token : ", GetApigeeToken())
	req.Header.Add("Authorization", "Bearer "+GetApigeeToken())
	return req, nil
}

func getHttpClient() (client *http.Client, err error) {

	if GetProxyURL() != "" {
		if pUrl, err := url.Parse(GetProxyURL()); err != nil {
			client = &http.Client{
				Transport: &http.Transport{
					Proxy: http.ProxyURL(pUrl),
				},
			}
		} else {
			return nil, err
		}
	} else {
		client = &http.Client{}
	}
	return client, nil
}

func handleResponse(print bool, resp *http.Response) (respBody []byte, err error) {

	if resp != nil {
		defer resp.Body.Close()
	}

	if resp == nil {
		clilog.Error.Println("error in response: Response was null")
		return nil, nil
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return nil, err
	} else if resp.StatusCode > 399 {
		clilog.Error.Printf("status code %d, error in response: %s\n", resp.StatusCode, string(respBody))
		return nil, errors.New("error in response")
	}
	if print {
		return respBody, PrettyPrint(respBody)
	}

	return respBody, nil
}
