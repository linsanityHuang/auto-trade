package main

import (
	"auto-trade/bigone"
	"auto-trade/global"
	"auto-trade/pkg/setting"
	"auto-trade/worker"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	btcUSDT = "BTC-USDT"
	ask     = "ASK"
	bid     = "BID"
	limit   = "LIMIT"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	setupHTTPClient()
}

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v\n", err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	tradeCmd := flag.NewFlagSet("trade", flag.ExitOnError)
	priceSpace := tradeCmd.Float64("priceSpace", 0.01, "price space")
	orderNum := tradeCmd.Int("orderNum", 0, "order number")

	orderCmd := flag.NewFlagSet("order", flag.ExitOnError)

	balanceCmd := flag.NewFlagSet("balance", flag.ExitOnError)
	asset := balanceCmd.String("asset", "BTC", "asset symbol")

	priceCmd := flag.NewFlagSet("price", flag.ExitOnError)
	assetPair := priceCmd.String("assetPair", "BTC-USDT", "asset pair")

	if len(os.Args) < 2 {
		fmt.Println("expected 'auto-trade', 'order', 'balance' or 'price' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "trade":
		tradeCmd.Parse(os.Args[2:])
		log.Printf("trade, price space: %.2f, order number: %d\n", *priceSpace, *orderNum)
		worker.StartTrade(float32(*priceSpace), *orderNum)
	case "order":
		orderCmd.Parse(os.Args[2:])
		bigone.ShowOrdes()
	case "balance":
		balanceCmd.Parse(os.Args[2:])
		r := bigone.SpotBalance(*asset)
		fmt.Printf("spot balance %s: %.2f\n", *asset, r)
	case "price":
		priceCmd.Parse(os.Args[2:])
		bigone.ShowAssetPrice(*assetPair)
	}
}

func setupSetting() error {
	setting, err := setting.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("BigOne", &global.BigOneSetting)
	if err != nil {
		return err
	}

	return nil
}

func setupHTTPClient() {
	global.HTTPClient = resty.New().
		SetHeader("Content-Type", "application/json").
		SetHostURL(global.BigOneSetting.BASEAPI).
		SetTimeout(2 * time.Second)
}
