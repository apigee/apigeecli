package crtapis

import (
	"../../shared"
	"github.com/spf13/cobra"
	"net/url"
	"path"
	"os"
	"archive/zip"
	"strings"
	"errors"
)

var Cmd = &cobra.Command{
	Use:   "create", 
	Short: "Creates an API proxy in an Apigee Org",
	Long:  "Creates an API proxy in an Apigee Org",
	Run: func(cmd *cobra.Command, args []string) {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")

		if proxy != "" {
			q := u.Query()			
			q.Set("name", name)
			q.Set("action", "import")	
			u.RawQuery = q.Encode()
			file, err := readProxyBundle()
			if err != nil {
				shared.PostHttpOctet(u.String(), file)
			}
			defer file.Close()					
		} else {
			proxyName := "{\"name\":\"" + name + "\"}"
			shared.PostHttpClient(u.String(), proxyName)	
		}
	},
}

var name, proxy string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "email of the developer")
	Cmd.Flags().StringVarP(&proxy, "proxy", "p",
		"", "API Proxy Bundle path")

	Cmd.MarkFlagRequired("name")
}

func readProxyBundle() (*os.File, error){

	if !strings.HasSuffix(proxy, ".zip") {
		shared.Error.Fatalln("Proxy bundle must be a zip file")
		return nil, errors.New("source must be a zipfile")
	}

	file, err := os.Open(proxy)

	if err != nil {
		shared.Error.Fatalln("Cannot open/read API Proxy Bundle: ", err)
		return nil, err
	}

	fi, err := file.Stat()
	_, err = zip.NewReader(file, fi.Size())

	if err != nil {
		shared.Error.Fatalln("Invalid API Proxy Bundle: ", err)
		return nil, err
	}

	return file, nil
}