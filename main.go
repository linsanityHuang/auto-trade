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

	"github.com/go-redis/redis"
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

	setupRedisClient()
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
	list := orderCmd.Bool("list", false, "pending order list")
	detail := orderCmd.Int64("show", 0, "show order detail by id")
	cancel := orderCmd.Int64("cancel", 0, "cancel order by id")

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
		log.Printf("trade, price space: %.5f, order number: %d\n", *priceSpace, *orderNum)
		worker.StartTrade(float32(*priceSpace), *orderNum)
	case "order":
		orderCmd.Parse(os.Args[2:])
		if *list {
			bigone.ShowOrdes()
		}

		if *detail > 0 {
			bigone.ShowOrde(*detail)
		}

		if *cancel > 0 {
			bigone.CancelOrder(*cancel)
		}
	case "balance":
		balanceCmd.Parse(os.Args[2:])
		r := bigone.SpotBalance(*asset)
		fmt.Printf("spot balance %s: %.5f\n", *asset, r)
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

	err = setting.ReadSection("Redis", &global.RedisSetting)
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

func setupRedisClient() {
	addr := fmt.Sprintf("%s:%d", global.RedisSetting.Host, global.RedisSetting.Port)
	global.RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       global.RedisSetting.DB,
		Password: "",
	})

	if err := global.RedisClient.Ping().Err(); err != nil {
		log.Fatalf("redis setup failed: %v\n", err)
	}
}
