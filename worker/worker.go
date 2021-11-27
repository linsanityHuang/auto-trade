package worker

import (
	"auto-trade/bigone"
	"auto-trade/global"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

func StartTrade(priceSpace float32, orderNums int) {
	// 读取现货账户USDT余额
	// totalUSDT := bigone.SpotUSDTBalance()
	var totalUSDT float32 = 20

	log.Printf("当前挂单使用的总金额为 %.2f USDT\n", totalUSDT)

	// 获取交易对当前价格
	ticker, err := bigone.ReadTicker(bigone.BTCUSDT)
	if err != nil {
		log.Fatalf("read ticker err: %v\n", err)
	}

	log.Printf("当前 %s 交易对的价格为 %s\n", bigone.BTCUSDT, ticker.Ask.Price)

	currentPrice, err := decimal.NewFromString(ticker.Ask.Price)
	p0, _ := currentPrice.BigFloat().Float32()

	if orderNums > 0 {
		// 生成网格买单
		buyOrders := bigone.CreateGridOrder(p0, priceSpace, totalUSDT, orderNums)

		_, err = bigone.CreateMultiOrder(buyOrders)
		if err != nil {
			log.Fatalf("create multi order err: %v\n", err)
		}
	}

	go Worker(priceSpace)

	timeTicker := time.NewTicker(time.Second)

	for {
		<-timeTicker.C
	}
}

func Worker(priceSpace float32) {
	for {

		orders, err := bigone.ReadOrders(bigone.BTCUSDT, bigone.PENDING)
		if err != nil {
			log.Printf("read exist order failed: %v\n", err)
		}

		for _, order := range orders {
			if order.State != "FILLED" {
				continue
			}

			log.Printf("当前订单, ID: %d, AssetPairName: %s, Side: %s, State: %s, Price: %s, Amount: %s\n",
				order.ID, order.AssetPairName, order.Side, order.State, order.Price, order.Amount)

			// 获取订单的成交价格
			orderPrice, _ := decimal.NewFromString(order.Price)

			if order.Side == bigone.BID.String() {
				sellPrice := orderPrice.Mul(decimal.NewFromFloat32(1 + priceSpace))

				_, err := bigone.CreateOrder("BTC-USDT", bigone.ASK.String(), sellPrice.String(), order.Amount, bigone.LIMIT.String())

				if err != nil {
					log.Printf("创建卖单失败 %v\n", err)
					continue
				}

				log.Printf("买单 > buyPrice: %s, amount: %s 对应的卖单已创建 > sellPrice: %s\n", order.Price, order.Amount, sellPrice.String())

			} else if order.Side == bigone.ASK.String() {
				buyPrice := orderPrice.Mul(decimal.NewFromFloat32(1 - priceSpace))

				_, err := bigone.CreateOrder("BTC-USDT", bigone.BID.String(), buyPrice.String(), order.Amount, bigone.LIMIT.String())

				if err != nil {
					log.Printf("创建买单失败 %v\n", err)
					continue
				}

				log.Printf("卖单 > sellPrice: %s, amount: %s 对应的卖单已创建 > sellPrice: %s\n", order.Price, order.Amount, buyPrice.String())
			}

		}

		time.Sleep(global.Timeout)
	}
}
