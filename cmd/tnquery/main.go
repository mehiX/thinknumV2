package main

import (
	"flag"
	"fmt"

	thinknum "github.com/mehiX/thinknumV2"
	"github.com/mehiX/thinknumV2/internal/query"
)

var (
	cfg = flag.String("c", "config.json", "Configuration file")
)

func main() {
	flag.Parse()

	fmt.Printf("Using configuration from %s\n", *cfg)

	client, err := thinknum.NewClient(*cfg)
	if err != nil {
		panic(err)
	}

	srch := client.Searches[0]

	hostname, version, token := client.Hostname, client.Version, client.Token
	datasetID := srch.DatasetID
	pageSize := client.PageSize
	ri, err := query.RunSearch(hostname, version, token, datasetID, pageSize, srch.Request)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Fields: %d\tRows: %d\n", len(ri.Fields), len(ri.Rows))

}
