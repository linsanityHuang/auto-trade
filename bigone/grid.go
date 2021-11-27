package bigone

import (
	"log"

	"github.com/shopspring/decimal"
)

// CreateGridOrder create grid order
func CreateGridOrder(p0, width, totalUSDT float32, orderNums int) (buyOrders []*Order) {
	// 批量创建订单数量在10个以内
	// buyOrders, sellOrders := []*Order{}, []*Order{}

	// 计算预期盈利
	diff, profit := ExpectedProfit(p0, width, totalUSDT, orderNums)
	log.Printf("此次交易预期收益: %.2fUSDT, 收益率: %.2f%%\n", diff, profit*100)

	startPrice := p0

	for i := 0; i < orderNums; i++ {
		buyOrders = append(buyOrders, &Order{
			AssetPairName: "BTC-USDT",
			Price:         decimal.NewFromFloat32(startPrice).String(),
			Side:          BID.String(),
			Type:          LIMIT.String(),
		})

		startPrice *= (1 - width)
	}

	for _, bOrder := range buyOrders {

		bPrice, _ := decimal.NewFromString(bOrder.Price)

		bOrder.Amount = decimal.NewFromFloat32(totalUSDT / float32(len(buyOrders))).Div(bPrice).String()
	}

	return
}

// ExpectedProfit 计算预期盈利
func ExpectedProfit(p0, width, totalUSDT float32, orderNum int) (float32, float32) {

	var afterUSDT float32

	startPrice := p0

	for i := 0; i < orderNum; i++ {
		// 买入价格
		startPrice *= (1 - width)

		// 买入BTC数量
		buyAmount := totalUSDT / float32(orderNum) / startPrice

		// 卖出获得USDT数量
		afterUSDT += (buyAmount * startPrice * (1 + width))

		log.Printf("buyPrice: %.2f, buyAmount: %.6f, sellPrice: %.2f", startPrice, buyAmount, startPrice*(1+width))
	}

	diff := afterUSDT - totalUSDT

	return diff, diff / totalUSDT
}
