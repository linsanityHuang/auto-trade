package bigone

import (
	"context"
	"fmt"
)

// Candle candle
type Candle struct {
	Close  string     `json:"close,omitempty"`
	High   PriceLevel `json:"high,omitempty"`
	Low    PriceLevel `json:"low,omitempty"`
	Open   string     `json:"open,omitempty"`
	Time   string     `json:"time,omitempty"`
	Volume string     `json:"volume,omitempty"`
}

// ReadCandles Candles of a asset pair
func ReadCandles(assetPairName, period string) ([]*Candle, error) {

	uri := fmt.Sprintf("/asset_pairs/%v/candles?period=%s", assetPairName, period)

	resp, err := HTTPRequest(context.Background()).Get(uri)

	if err != nil {
		return nil, err
	}

	candles := []*Candle{}

	if err := UnmarshalResponse(resp, &candles); err != nil {
		return nil, err
	}

	return candles, nil
}
