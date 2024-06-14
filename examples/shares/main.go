package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	investgo "github.com/floatdrop/tbank-invest-go-sdk/sdk"
)

func usage() {
	log.Printf("Usage: shares [-e endpoint] -t token\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var endpoint = flag.String("e", "sandbox-invest-public-api.tinkoff.ru:443", "TBank API endpoint")
	var token = flag.String("t", "", "TBank invest token (from https://www.tinkoff.ru/invest/settings/)")
	var showHelp = flag.Bool("h", false, "Show help message")
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	client, err := investgo.NewClient(*endpoint, *token)
	if err != nil {
		log.Fatal("Could not connect: %w", err)
	}
	defer func() { _ = client.Close() }()

	shares, err := client.AllShares(context.Background())
	if err != nil {
		log.Fatal("Failed to get shares: %w", err)
	}

	for _, share := range shares {
		fmt.Println(share.Uid, share.Name)
	}
}
