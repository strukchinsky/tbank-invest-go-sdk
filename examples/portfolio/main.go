package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/strukchinsky/tbank-invest-go-sdk"
	investgo "github.com/strukchinsky/tbank-invest-go-sdk/sdk"
)

func usage() {
	log.Printf("Usage: portfolio [-e endpoint] -t token [ACCOUNT_ID]\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

func main() {
	var endpoint = flag.String("e", "invest-public-api.tinkoff.ru:443", "TBank API endpoint")
	var token = flag.String("t", "", "TBank invest token (from https://www.tinkoff.ru/invest/settings/)")

	var showHelp = flag.Bool("h", false, "Show help message")
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	client, err := investgo.NewClient(*endpoint, *token)
	if err != nil {
		log.Fatal("Failed to get accounts: ", err)
	}
	defer func() { _ = client.Close() }()

	args := flag.Args()

	if len(args) != 1 {
		log.Println("Fetching accounts list...")
		accounts, err := client.GetAccounts(context.Background())
		if err != nil {
			log.Fatal("Could not get accounts: ", err)
		}

		fmt.Println("Use Id as first argument to select account:")
		for _, account := range accounts {
			fmt.Printf("Id=%s Name=%s\n", account.Id, account.Name)
		}

		os.Exit(0)
	}

	log.Println("Fetching positions for portfolio...")
	portfolio, err := client.GetPortfolio(context.Background(), args[0], pb.PortfolioRequest_RUB)
	if err != nil {
		log.Fatal("Failed to get portfolio for account: ", err)
	}

	fmt.Println("Portfolio structure:")
	for _, position := range portfolio.Positions {
		instrument, err := client.InstrumentBy(context.Background(), investgo.ByUid(position.InstrumentUid))
		if err != nil {
			log.Fatal("Failed to get instrument info: ", err)
		}

		positionValue := investgo.QuotationToFloat(position.Quantity) * investgo.MoneyValueToFloat(position.CurrentPrice)
		percentage := positionValue / investgo.MoneyValueToFloat(portfolio.TotalAmountPortfolio)

		fmt.Printf("(%05.2f%%) %s\n", 100*percentage, instrument.Name)
	}
}
