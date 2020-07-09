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
	"os"

	"github.com/srinandan/apigeecli/clilog"
)

//PostHttpOctet method is used to send resources, proxy bundles, shared flows etc.
func PostHttpOctet(print bool, url string, proxyName string) (respBody []byte, err error) {
	file, _ := os.Open(proxyName)
	defer file.Close()

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
	client := &http.Client{}

	clilog.Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	clilog.Info.Println("Setting token : ", GetApigeeToken())
	req.Header.Add("Authorization", "Bearer "+GetApigeeToken())
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	defer resp.Body.Close()
	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return nil, err
	} else if resp.StatusCode != 200 {
		clilog.Error.Println("error in response: ", string(respBody))
		return nil, errors.New("error in response")
	}
	if print {
		return respBody, PrettyPrint(respBody)
	}

	return respBody, nil
}

//DownloadResource method is used to download resources, proxy bundles, sharedflows
func DownloadResource(url string, name string, resType string) error {
	var filename string

	if resType == ".zip" {
		filename = name + ".zip"
	} else {
		filename = name
	}

	out, err := os.Create(filename)
	if err != nil {
		clilog.Error.Println("error creating file: ", err)
		return err
	}
	defer out.Close()

	client := &http.Client{}

	clilog.Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return err
	}

	clilog.Info.Println("Setting token : ", GetApigeeToken())
	req.Header.Add("Authorization", "Bearer "+GetApigeeToken())

	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return err
	} else if resp.StatusCode > 299 {
		clilog.Error.Println("error in response: ", resp.Body)
		return errors.New("error in response")
	}
	defer resp.Body.Close()
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

	client := &http.Client{}
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
		contentType = params[2]
	default:
		return nil, errors.New("unsupported method")
	}

	if err != nil {
		clilog.Error.Println("error in client: ", err)
		return nil, err
	}

	if GetApigeeToken() == "" {
		if err = SetAccessToken(); err != nil {
			return nil, err
		}
	}

	clilog.Info.Println("Setting token : ", GetApigeeToken())
	req.Header.Add("Authorization", "Bearer "+GetApigeeToken())
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)

	if err != nil {
		clilog.Error.Println("error connecting: ", err)
		return nil, err
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		clilog.Error.Println("error in response: ", err)
		return nil, err
	} else if resp.StatusCode > 299 {
		clilog.Error.Println("error in response: ", string(respBody))
		return nil, errors.New("error in response")
	}
	if print {
		return respBody, PrettyPrint(respBody)
	}
	return respBody, nil
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
	} else {
		return nil, errors.New("unsupported method")
	}
	return req, err
}
