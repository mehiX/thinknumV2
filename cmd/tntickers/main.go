package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	thinknum "github.com/mehiX/thinknumV2"
)

var cfg = flag.String("c", "config.json", "Configuration file to use")
var dataset = flag.String("d", "", "Dataset ID")

func main() {
	flag.Parse()

	if *dataset == "" {
		fmt.Println("No dataset provided")
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("Using configuration from: %s\n", *cfg)

	client, err := thinknum.NewClient(*cfg)
	if err != nil {
		panic(err)
	}

	tickers, err := client.Tickers(*dataset)
	if err != nil {
		panic(err)
	}

	log.Println("Retrieved", len(tickers))
	//for _, t := range tickers {
	//fmt.Printf("%20s %s\n", t.ID, t.DisplayName)
	//}
}
