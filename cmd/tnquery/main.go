package main

import (
	"flag"
	"fmt"

	thinknum "github.com/mehiX/thinknumV2"
)

var (
	cfg = flag.String("c", "config.json", "Configuration file")
)

func main() {
	flag.Parse()

	fmt.Printf("Using configuration from %s\n", *cfg)

	client, err := thinknum.NewClientFromJSON(*cfg)
	if err != nil {
		panic(err)
	}

	for ri := range client.RunAll() {
		if ri.Error != nil {
			fmt.Printf("Error: %v\n", ri.Error)
		}
		fmt.Printf("Output to: %s\n", ri.Search.OutputFile)
		fmt.Printf("Fields: %d\tRows: %d/%d\tPages: %d\n", len(ri.Data.Fields), len(ri.Data.Rows), ri.Data.Total, ri.Data.Pages)
	}

}
