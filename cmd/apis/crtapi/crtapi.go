package crtapis

import (
	"../../shared"
	"archive/zip"
	"errors"
	"github.com/spf13/cobra"
	"net/url"
	"os"
	"path"
	"strings"
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
			err := readProxyBundle()
			if err == nil {
				_ = shared.PostHttpOctet(u.String(), proxy)
			}
		} else {
			proxyName := "{\"name\":\"" + name + "\"}"
			_ = shared.HttpClient(u.String(), proxyName)
		}
	},
}

var name, proxy string

func init() {

	Cmd.Flags().StringVarP(&name, "name", "n",
		"", "API Proxy name")
	Cmd.Flags().StringVarP(&proxy, "proxy", "p",
		"", "API Proxy Bundle path")

	_ = Cmd.MarkFlagRequired("name")
}

func readProxyBundle() error {

	if !strings.HasSuffix(proxy, ".zip") {
		shared.Error.Fatalln("Proxy bundle must be a zip file")
		return errors.New("source must be a zipfile")
	}

	file, err := os.Open(proxy)

	if err != nil {
		shared.Error.Fatalln("Cannot open/read API Proxy Bundle: ", err)
		return err
	}

	fi, err := file.Stat()
	if err != nil {
		shared.Error.Fatalln("Error accessing file: ", err)
		return err
	}
	_, err = zip.NewReader(file, fi.Size())

	if err != nil {
		shared.Error.Fatalln("Invalid API Proxy Bundle: ", err)
		return err
	}

	defer file.Close()
	return nil
}
