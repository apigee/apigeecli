package expapis

import (
	"encoding/json"
	"net/url"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

type proxies struct {
	Proxies []proxy `json:"proxies,omitempty"`
}

type proxy struct {
	Name     string   `json:"name,omitempty"`
	Revision []string `json:"revision,omitempty"`
}

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "export API proxy bundles from an org",
	Long:  "export API proxy bundles from an org",
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return exportProxies()
	},
}

var conn int

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")
}

func batch(entities []proxy, entityType string, pwg *sync.WaitGroup) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		//download only the last revision
		lastRevision := len(entity.Revision)
		go shared.FetchAsyncBundle(entityType, entity.Name, entity.Revision[lastRevision-1], &bwg)
	}
	bwg.Wait()
}

func exportProxies() error {

	//parent workgroup
	var pwg sync.WaitGroup
	const entityType = "apis"

	u, _ := url.Parse(shared.BaseURL)
	q := u.Query()
	q.Set("includeRevisions", "true")
	u.RawQuery = q.Encode()
	u.Path = path.Join(u.Path, shared.RootArgs.Org, entityType)

	//don't print to sysout
	respBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	var entities = proxies{}
	err = json.Unmarshal(respBody, &entities)
	if err != nil {
		return err
	}

	numEntities := len(entities.Proxies)
	shared.Info.Printf("Found %d API Proxies in the org\n", numEntities)
	shared.Info.Printf("Exporting bundles with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than products
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Exporting batch %d of proxies\n", (i + 1))
		go batch(entities.Proxies[start:end], entityType, &pwg)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Exporting remaining %d proxies\n", remaining)
		go batch(entities.Proxies[start:numEntities], entityType, &pwg)
		pwg.Wait()
	}

	return nil
}
