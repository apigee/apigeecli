package expts

import (
	"encoding/json"
	"net/url"
	"path"
	"sync"

	"github.com/spf13/cobra"
	"github.com/srinandan/apigeecli/cmd/shared"
)

////Cmd to export target servers
var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export target servers to a file",
	Long:  "Export target servers to a file",
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		const exportFileName = "targetservers.json"
		err = exportTargetServers()
		if err != nil {
			return
		}
		return shared.WriteByteArrayToFile(exportFileName, false, nil)
	},
}

var conn int

func init() {

	Cmd.Flags().IntVarP(&conn, "conn", "c",
		4, "Number of connections")

}

//batch created a batch of targetservers to query
func batch(entities []string, entityType string, pwg *sync.WaitGroup, mu *sync.Mutex) {

	defer pwg.Done()
	//batch workgroup
	var bwg sync.WaitGroup

	bwg.Add(len(entities))

	for _, entity := range entities {
		u, _ := url.Parse(shared.BaseURL)
		u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, entityType, entity)
		go shared.GetAsyncEntity(u.String(), &bwg, mu)
	}
	bwg.Wait()
}

func exportTargetServers() error {

	//parent workgroup
	var pwg sync.WaitGroup
	var mu sync.Mutex
	const entityType = "targetservers"

	u, _ := url.Parse(shared.BaseURL)
	u.Path = path.Join(u.Path, shared.RootArgs.Org, "environments", shared.RootArgs.Env, entityType)
	//don't print to sysout
	respBody, err := shared.HttpClient(false, u.String())
	if err != nil {
		return err
	}

	var targetservers []string
	err = json.Unmarshal(respBody, &targetservers)
	if err != nil {
		return err
	}

	numEntities := len(targetservers)
	shared.Info.Printf("Found %d targetservers in the org\n", numEntities)
	shared.Info.Printf("Exporting targetservers with %d connections\n", conn)

	numOfLoops, remaining := numEntities/conn, numEntities%conn

	//ensure connections aren't greater than targetservers
	if conn > numEntities {
		conn = numEntities
	}

	start := 0

	for i, end := 0, 0; i < numOfLoops; i++ {
		pwg.Add(1)
		end = (i * conn) + conn
		shared.Info.Printf("Exporting batch %d of targetservers\n", (i + 1))
		go batch(targetservers[start:end], entityType, &pwg, &mu)
		start = end
		pwg.Wait()
	}

	if remaining > 0 {
		pwg.Add(1)
		shared.Info.Printf("Exporting remaining %d targetservers\n", remaining)
		go batch(targetservers[start:numEntities], entityType, &pwg, &mu)
		pwg.Wait()
	}

	return nil
}
