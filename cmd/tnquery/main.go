package main

import (
	"flag"
	"fmt"
	"sync"

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

	for ri := range run(client) {
		if ri.Error != nil {
			fmt.Printf("Error: %v\n", ri.Error)
		}
		fmt.Printf("Output to: %s\n", ri.OutputFile)
		fmt.Printf("Fields: %d\tRows: %d/%d\tPages: %d\n", len(ri.Data.Fields), len(ri.Data.Rows), ri.Data.Total, ri.Data.Pages)
	}

}

type runresult struct {
	Data       query.RowsItems
	Error      error
	OutputFile string
}

func run(client *thinknum.Client) <-chan runresult {

	resultsStream := make(chan runresult)

	go func() {
		defer close(resultsStream)

		searchesStream := make(chan thinknum.SearchDefinition)
		// generate work
		go func() {
			defer close(searchesStream)

			// from slice to channel
			for _, s := range client.Searches {
				// skip disabled seaches
				if !s.Disabled {
					searchesStream <- s
				} else {
					fmt.Printf("Skip disabled search: %s\n", s.Name)
				}
			}
		}()

		workers := client.Workers

		var wg sync.WaitGroup
		wg.Add(workers)

		// start idle workers
		for i := 0; i < workers; i++ {
			go runner(client, searchesStream, &wg, resultsStream)
		}

		wg.Wait()
	}()

	return resultsStream

}

func runner(client *thinknum.Client, searches <-chan thinknum.SearchDefinition, wg *sync.WaitGroup, results chan<- runresult) {
	defer wg.Done()

	for s := range searches {
		ri, err := client.RunSearch(s)
		results <- runresult{ri, err, s.OutputFile}
	}
}
