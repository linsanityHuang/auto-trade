package bigone

import (
	"context"
	"fmt"
)

// PriceLevel PriceLevel
type PriceLevel struct {
	Price      string `json:"price,omitempty"`
	Quantity   string `json:"quantity,omitempty"`
	OrderCount int    `json:"order_count,omitempty"`
}

// Ticker Ticker
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

// ReadTicker ReadTicker
func ReadTicker(assetPairName string) (*Ticker, error) {

	resp, err := HTTPRequest(context.Background()).Get(fmt.Sprintf("/asset_pairs/%v/ticker", assetPairName))
	if err != nil {
		return nil, err
	}

	ticker := &Ticker{}

	if err := UnmarshalResponse(resp, ticker); err != nil {
		return nil, err
	}

	return ticker, nil
}

// ReadTickers ReadTickers
func ReadTickers(assetPairName string) ([]*Ticker, error) {

	resp, err := HTTPRequest(context.Background()).Get(fmt.Sprintf("/asset_pairs/tickers?pair_names=%s", assetPairName))

	if err != nil {
		return nil, err
	}

	tickers := []*Ticker{}

	if err := UnmarshalResponse(resp, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
