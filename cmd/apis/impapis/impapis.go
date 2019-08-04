package impapis

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
	"github.com/srinandan/apigeecli/cmd/types"
)

var Cmd = &cobra.Command{
	Use:   "import",
	Short: "Import a folder containing an API proxy bundles",
	Long:  "Import a folder containing an API proxy bundles",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		u, _ := url.Parse(shared.BaseURL)
		return createAPIs(u)
	},
}

var folder string
var conn int

func init() {

	Cmd.Flags().StringVarP(&folder, "folder", "f",
		"", "folder containing API proxy bundles")
	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

	_ = Cmd.MarkFlagRequired("folder")
}

func createAsyncAPI(u *url.URL, bundlePath string, wg *sync.WaitGroup, errChan chan<- *types.ImportError) {

	defer wg.Done()
	_, fileName := filepath.Split(bundlePath)
	name := strings.Split(fileName, ".")

	q := u.Query()
	q.Set("name", name[0])
	q.Set("action", "import")
	u.RawQuery = q.Encode()
	err := shared.ReadBundle(bundlePath)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	_, err = shared.PostHttpOctet(true, u.String(), bundlePath)
	if err != nil {
		errChan <- &types.ImportError{Err: err}
		return
	}

	errChan <- &types.ImportError{Err: nil}
}

func createAPIs(u *url.URL) error {

	var errChan = make(chan *types.ImportError)
	var wg sync.WaitGroup
	var proxyBundles []string

	u.Path = path.Join(u.Path, shared.RootArgs.Org, "apis")

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".zip" {
			return nil
		}
		proxyBundles = append(proxyBundles, path)
		return nil
	})

	if err != nil {
		return err
	}

	numAPIs := len(proxyBundles)
	shared.Info.Printf("Found %d bundles in the folder\n", numAPIs)
	shared.Info.Printf("Create proxies with %d connections\n", conn)

	if numAPIs < conn {
		wg.Add(numAPIs)
		for i := 0; i < numAPIs; i++ {
			go createAsyncAPI(u, proxyBundles[i], &wg, errChan)
		}

		go func() {
			wg.Wait()
			close(errChan)
		}()

	} else {
		numOfLoops, remaining := numAPIs/conn, numAPIs%conn
		for i := 0; i < numOfLoops; i++ {
			shared.Info.Printf("Create %d batch of proxies\n", i)
			wg.Add(conn)
			for j := 0; j < conn; j++ {
				go createAsyncAPI(u, proxyBundles[j+(i*conn)], &wg, errChan)
			}
			go func() {
				wg.Wait()
			}()
		}

		wg.Add(remaining)
		shared.Info.Printf("Create remaining %d proxies\n", remaining)
		for i := (numAPIs - remaining); i < numAPIs; i++ {
			go createAsyncAPI(u, proxyBundles[i], &wg, errChan)
		}
		go func() {
			wg.Wait()
			close(errChan)
		}()
	}

	//print any errors and return an err
	var errs = false
	for errAPIs := range errChan {
		if errAPIs.Err != nil {
			shared.Error.Fatalln(errAPIs.Err)
			errs = true
		}
	}

	if errs {
		return fmt.Errorf("problem creating one of more products")
	}
	return nil
}
