package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	investapi "github.com/floatdrop/tbank-invest-go-sdk"
	investgo "github.com/floatdrop/tbank-invest-go-sdk/sdk"
)

func usage() {
	log.Printf("Usage: candles [-e endpoint] -t token INSTRUMENT_UID\n")
	flag.PrintDefaults()
}

func showUsageAndExit(exitcode int) {
	usage()
	os.Exit(exitcode)
}

type DateValue struct {
	Date *time.Time
}

const dateLayout = "2006-01-02 15:04"

func (v DateValue) String() string {
	if v.Date != nil {
		return v.Date.Format(dateLayout)
	}
	return ""
}

func (v DateValue) Set(str string) error {
	if t, err := time.Parse(dateLayout, str); err != nil {
		return err
	} else {
		*v.Date = t
	}
	return nil
}

func main() {
	var endpoint = flag.String("e", "sandbox-invest-public-api.tinkoff.ru:443", "TBank API endpoint")
	var token = flag.String("t", "", "TBank invest token (from https://www.tinkoff.ru/invest/settings/)")

	var from = time.Now().Add(-1 * time.Hour * 24)
	flag.Var(&DateValue{&from}, "from", "")

	var to = time.Now()
	flag.Var(&DateValue{&to}, "to", "")

	var showHelp = flag.Bool("h", false, "Show help message")
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		showUsageAndExit(0)
	}

	args := flag.Args()
	if len(args) != 1 {
		showUsageAndExit(1)
	}

	client, err := investgo.NewClient(*endpoint, *token)
	if err != nil {
		log.Fatal("Could not connect: ", err)
	}
	defer func() { _ = client.Close() }()

	fmt.Println("Fetching candles for", args[0])
	candles, err := client.GetCandles(
		context.Background(),
		args[0],
		from,
		to,
		investapi.CandleInterval_CANDLE_INTERVAL_1_MIN,
	)
	if err != nil {
		log.Fatal("Failed to get shares: ", err)
	}

	for _, candle := range candles {
		fmt.Printf("[%s] close=%d\n", candle.Time.AsTime().Format(dateLayout), candle.Close.Units)
	}
}
