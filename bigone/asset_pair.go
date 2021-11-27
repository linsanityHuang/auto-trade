package bigone

import (
	"context"
)

const (
	BTCUSDT string = "BTC-USDT"
	USDT    string = "USDT"
	BTC     string = "BTC"
)

// Asset asset
type Asset struct {
	ID     string `json:"id,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Name   string `json:"name,omitempty"`
}

// AssetPair Account represents the state of one asset pair.
type AssetPair struct {
	ID            string `json:"id,omitempty"`
	QuoteScale    int    `json:"quote_scale,omitempty"`
	QuoteAsset    Asset  `json:"quote_asset,omitempty"`
	Name          string `json:"name,omitempty"`
	BaseScale     int    `json:"base_scale,omitempty"`
	MinQuoteValue string `json:"min_quote_value,omitempty"`
	BaseAsset     Asset  `json:"base_asset,omitempty"`
}

// ReadAssetPairs All AssetPairs
func ReadAssetPairs() ([]*AssetPair, error) {
	resp, err := HTTPRequest(context.Background()).Get("/asset_pairs")

	if err != nil {
		return nil, err
	}

	assetPairs := []*AssetPair{}

	if err := UnmarshalResponse(resp, &assetPairs); err != nil {
		return nil, err
	}

	return assetPairs, nil
}
