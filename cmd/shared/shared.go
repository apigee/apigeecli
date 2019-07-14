package shared

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const BaseURL = "https://apigee.googleapis.com/v1/organizations/"

// Arguements is the base struct to hold all command arguments
type Arguments struct {
	Verbose        bool
	Org            string
	Env            string
	Token          string
	ServiceAccount string
}

var RootArgs = Arguments{}

//log levels, default is error
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var LogInfo = false

// Structure to hold OAuth response
type OAuthAccessToken struct {
	AccessToken string `json:"access_token,omitempty"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
}

//Init function initializes the logger objects
func Init() {

	var infoHandle = ioutil.Discard

	if LogInfo {
		infoHandle = os.Stdout
	}

	warningHandle := os.Stdout
	errorHandle := os.Stdout

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func GetHttpClient(url string, token string) error {
	client := &http.Client{}

	Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("GET", url, nil)

	Info.Println("Setting token : ", token)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("Error connecting:\n", err)
		return err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Error.Fatalln("Error in response:\n", err)
			return err
		} else if resp.StatusCode != 200 {
			Error.Fatalln("Error in response:\n", string(body))
			return errors.New("Error in response")
		} else {
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, body, "", "\t")
			if err != nil {
				Error.Fatalln("Error parsing response:\n", err)
				return err
			}
			fmt.Println(string(prettyJSON.Bytes()))
			return nil
		}
	}
}

func PostHttpClient(url string, token string, payload string) error {
	client := &http.Client{}

	Info.Println("Connecting to : ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))

	Info.Println("Setting token : ", token)
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("Error connecting:\n", err)
		return err
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			Error.Fatalln("Error in response:\n", err)
			return err
		} else if resp.StatusCode != 200 {
			Error.Fatalln("Error in response:\n", string(body))
			return errors.New("Error in response")
		} else {
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, body, "", "\t")
			if err != nil {
				Error.Fatalln("Error parsing response:\n", err)
				return err
			}
			fmt.Println(string(prettyJSON.Bytes()))
			return nil
		}
	}
}

func getPrivateKey() (interface{}, error) {
	pemPrivateKey := fmt.Sprintf("%v", viper.Get("private_key"))
	block, _ := pem.Decode([]byte(pemPrivateKey))
	privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Error.Fatalln("Error parsing Private Key:\n", err)
		return nil, err
	} else {
		return privKey, nil
	}
}

func generateJWT() (string, error) {

	const aud = "https://www.googleapis.com/oauth2/v4/token"
	const scope = "https://www.googleapis.com/auth/cloud-platform"

	privKey, err := getPrivateKey()

	if err != nil {
		return "", err
	}

	now := time.Now()
	token := jwt.New()

	token.Set(jwt.AudienceKey, aud)
	token.Set(jwt.IssuerKey, viper.Get("client_email"))
	token.Set("scope", scope)
	token.Set(jwt.IssuedAtKey, now.Unix())
	token.Set(jwt.ExpirationKey, now.Unix())

	payload, err := token.Sign(jwa.RS256, privKey)
	if err != nil {
		Error.Fatalln("Error parsing Private Key:\n", err)
		return "", err
	} else {
		Info.Println("jwt token : ", string(payload))
		return string(payload), nil
	}
}

func GenerateAccessToken() (string, error) {

	const token_endpoint = "https://www.googleapis.com/oauth2/v4/token"
	const grant_type = "urn:ietf:params:oauth:grant-type:jwt-bearer"

	token, err := generateJWT()

	if err != nil {
		return "", nil
	}

	form := url.Values{}
	form.Add("grant_type", grant_type)
	form.Add("assertion", token)

	client := &http.Client{}
	req, err := http.NewRequest("POST", token_endpoint, strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(req)

	if err != nil {
		Error.Fatalln("Failed to generate oauth token: \n", err)
		return "", err
	} else {
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			Error.Fatalln("Error in response: \n", string(bodyBytes))
			return "", errors.New("Error in response")
		} else {
			decoder := json.NewDecoder(resp.Body)
			accessToken := OAuthAccessToken{}
			if err := decoder.Decode(&accessToken); err != nil {
				Error.Fatalln("Error in response: \n", err)
				return "", errors.New("Error in response")
			} else {
				Info.Println("access token : ", accessToken)
				RootArgs.Token = accessToken.AccessToken
				return accessToken.AccessToken, nil
			}
		}
	}
}
