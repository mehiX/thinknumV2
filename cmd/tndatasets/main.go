package main

import (
	"flag"
	"fmt"
	"log"

	thinknum "github.com/mehiX/thinknumV2"
)

var (
	cfg      = flag.String("c", "config.json", "File to load configuration from")
	tickerID = flag.String("t", "", "TickerID to filter datasets by")
)

func main() {

	flag.Parse()

	fmt.Printf("Using configuration from: %s\n", *cfg)

	client, err := thinknum.NewClient(*cfg)
	if err != nil {
		log.Fatalln(err)
	}

	ds, err := client.Datasets(*tickerID)
	if err != nil {
		log.Fatalln(err)
	}

	for _, d := range ds {
		fmt.Printf("%-30s %s\n", d.ID, d.DisplayName)
	}
}
