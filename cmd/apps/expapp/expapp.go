package expapp

import (
	"encoding/json"
	"net/url"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

//Cmd to export apps
var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export API products to a file",
	Long:  "Export API products to a file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		const exportFileName = "apps.json"
		err = exportApps()
		if err != nil {
			return
		}
		return shared.WriteByteArrayToFile(exportFileName, false, nil)
	},
}

type apps struct {
	Apps []app `json:"app,omitempty"`
}

type app struct {
	AppID string `json:"appId,omitempty"`
}

var conn int

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

}

//batch created a batch of apps to query
func batch(entities []app, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, entityType, entity.AppID)
		go shared.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

func exportApps() error {

	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "apps"

	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, entityType)
	//don't print to sysout
	respBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	var entities = apps{}
	err = json.Unmarshal(respBody, &entities)
	if err != nil {
		return err
	}

	numEntities := len(entities.Apps)
	shared.Info.Printf("Found %d apps in the org\n", numEntities)
	shared.Info.Printf("Exporting apps with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than apps
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Exporting batch %d of apps\n", (i + 1))
		go batch(entities.Apps[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Exporting remaining %d apps\n", remaining)
		go batch(entities.Apps[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return nil
}
