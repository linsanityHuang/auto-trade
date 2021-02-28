package main

import (
	"auto-trade/bigone"
	"fmt"

	"github.com/shopspring/decimal"
)

// CreateGridOrder create grid order
func CreateGridOrder(p0, lowPrice, highPrice, width, totalUSDT float32) ([]*bigone.Order, []*bigone.Order) {
	buyOrders, sellOrders := []*bigone.Order{}, []*bigone.Order{}

	startPrice := p0

	for startPrice >= lowPrice {

		startPrice *= (1 - width)

		buyOrders = append(buyOrders, &bigone.Order{
			AssetPairName: "BTC-USDT",
			Price:         decimal.NewFromFloat32(startPrice).String(),
			Side:          "BID",
			Type:          "LIMIT",
			Operator:      "LTE",
		})
	}

	for _, bOrder := range buyOrders {

		bOrder.Amount = decimal.NewFromFloat32(totalUSDT / float32(len(buyOrders))).String()

		bPrice, _ := decimal.NewFromString(bOrder.Price)

		sellOrders = append(sellOrders, &bigone.Order{
			AssetPairName: "BTC-USDT",
			Price:         bPrice.Mul(decimal.NewFromFloat32(1 + width)).String(),
			Amount:        bOrder.Amount,
			Side:          "ASK",
			Type:          "LIMIT",
			Operator:      "GTE",
		})
	}

	return buyOrders, sellOrders
}

func main() {
	buyOrders, sellOrders := CreateGridOrder(10000, 5000, 15000, 0.1, 5000)

	for i, p := range buyOrders {
		fmt.Println(p)
		fmt.Println(sellOrders[i])
		fmt.Println("--------------------")
	}
}
