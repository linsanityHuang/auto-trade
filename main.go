package main

import (
	"auto-trade/bigone"
	"auto-trade/global"
	"auto-trade/pkg/setting"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

const (
	btcUSDT = "BTC-USDT"
	ask     = "ASK"
	bid     = "BID"
	limit   = "LIMIT"
)

func minQuoteValue() (string, error) {
	// 获取交易对最小下单数量
	assetPairs, err := bigone.ReadAssetPairs()
	if err != nil {
		return "", err
	}

	for _, pair := range assetPairs {
		if pair.Name == btcUSDT {
			fmt.Println(pair)
			return pair.MinQuoteValue, nil
		}
	}

	return "", nil
}

func readUSDTBalance() float32 {
	spotAccount, err := bigone.ReadSpotAccount("USDT")
	if err != nil {
		log.Fatalf("read spot account err: %v\n", err)
	}

	balance, err := decimal.NewFromString(spotAccount.Balance)
	if err != nil {
		log.Fatalf("balance to decimal err: %v\n", err)
	}

	b, _ := balance.BigFloat().Float32()

	if b < 10 {
		log.Fatalf("spot account usdt balance too small: %v\n", spotAccount.Balance)
	}

	return b
}

func gridTradeRun(width float32) {
	// 读取现货账户USDT余额
	// totalUSDT := readUSDTBalance()
	var totalUSDT float32 = 10

	// 获取交易对当前价格
	ticker, err := bigone.ReadTicker(btcUSDT)
	if err != nil {
		log.Fatalf("read ticker err: %v\n", err)
	}

	currentPrice, err := decimal.NewFromString(ticker.Bid.Price)
	p0, _ := currentPrice.BigFloat().Float32()

	// 生成网格买单
	buyOrders := bigone.CreateGridOrder(p0, width, totalUSDT, 5)

	pendingOrders, err := bigone.CreateMultiOrder(buyOrders)
	if err != nil {
		log.Fatalf("create multi order err: %v\n", err)
	}

	// 记录低价买单是否已经创建高价卖单
	exist := map[int64]bool{}

	for {

		for _, o := range pendingOrders {

			// 查询买单是否成交
			order, err := bigone.ReadOrder(o.ID)
			if err != nil {
				log.Printf("查询订单状态 err: %v\n", err)
				continue
			}

			// 跳过未成交买单
			if order.State != "FILLED" {
				continue
			}

			if _, ok := exist[o.ID]; ok {
				continue
			}

			// 创建网格卖单
			buyPrice, _ := decimal.NewFromString(order.Price)
			sellPrice := buyPrice.Mul(decimal.NewFromFloat32(1 + width))

			_, err = bigone.CreateOrder("BTC-USDT", "ASK", sellPrice.String(), order.Amount, "LIMIT")

			exist[o.ID] = true

			if err != nil {
				log.Printf("创建卖单失败 err: %v\n", err)
				continue
			}
		}

		// 如果所有买单都已成交, 且对应卖单已创建(无论成功与否)
		// 退出循环
		if len(exist) == len(buyOrders) {
			break
		}

		time.Sleep(5 * time.Second)
	}

	log.Println("grid trade exit.")
}

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
		panic(err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// gridTradeRun(0.02)

	fmt.Println(readUSDTBalance())
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
