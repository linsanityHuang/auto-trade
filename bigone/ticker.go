package bigone

import (
	"context"
	"fmt"
	"strings"
)

// PriceLevel PriceLevel
type PriceLevel struct {
	Price      string `json:"price,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
	OrderCount int    `json:"order_count,omitempty"`
}

// Ticker 是一个资产对的当前状态，有24小时的交易数据
type Ticker struct {
	AssetPairName string     `json:"asset_pair_name,omitempty"`
	Bid           PriceLevel `json:"bid,omitempty"`
	Ask           PriceLevel `json:"ask,omitempty"`
	Open          string     `json:"open,omitempty"`
	Close         string     `json:"close,omitempty"`
	High          string     `json:"high,omitempty"`
	Low           string     `json:"low,omitempty"`
	Volume        string     `json:"volume,omitempty"`
	DailyChange   string     `json:"daily_change,omitempty"`
}

// ReadTicker Ticker of one asset pair
func ReadTicker(assetPairName string) (*Ticker, error) {

	path := fmt.Sprintf("/asset_pairs/%v/ticker", strings.ToUpper(assetPairName))

	resp, err := HTTPRequest(context.Background()).Get(path)
	if err != nil {
		return nil, err
	}

	ticker := &Ticker{}

	if err := UnmarshalResponse(resp, ticker); err != nil {
		return nil, err
	}

	return ticker, nil
}

// ReadTickers Ticker of multiple asset pairs
func ReadTickers(pairNames []string) ([]*Ticker, error) {
	// pairNames BTC-USDT,PCX-BTC,GXC-USDT
	strings.Join(pairNames, ",")

	path := fmt.Sprintf("/asset_pairs/tickers?pair_names=%s", strings.Join(pairNames, ","))

	resp, err := HTTPRequest(context.Background()).Get(path)

	if err != nil {
		return nil, err
	}

	tickers := []*Ticker{}

	if err := UnmarshalResponse(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}

func ShowAssetPrice(assetPair string) {
	pair := strings.ToUpper(assetPair)
	// 获取交易对当前价格
	ticker, err := ReadTicker(pair)
	if err != nil {
		fmt.Printf("read ticker err: %v\n", err)
	}

	fmt.Printf("当前 %s 交易对的价格为 %s\n", pair, ticker.Ask.Price)
}
