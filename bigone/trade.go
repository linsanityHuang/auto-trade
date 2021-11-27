package bigone

import (
	"context"
	"fmt"
	"strings"
)

// Trade 已成交的交易订单
type Trade struct {
	ID            int64  `json:"id,omitempty"`              // id of trade
	AssetPairName string `json:"asset_pair_name,omitempty"` // asset pair name
	Price         string `json:"price,omitempty"`           // deal price
	Amount        string `json:"amount,omitempty"`          // amount
	TakerSide     string `json:"taker_side,omitempty"`      // order side, one of "ASK/BID"
	Side          string `json:"side,omitempty"`            // viewer side, one of "ASK/BID/SELF_TRADING"
	MakerOrderID  int64  `json:"maker_order_id,omitempty"`  // maker order id, null if taker_side is equal to side
	TakerOrderID  int64  `json:"taker_order_id,omitempty"`  // taker order id, null if taker_side is not equal to side
	MakerFee      string `json:"maker_fee,omitempty"`       // maker fee, null if taker_side is equal to side
	TakerFee      string `json:"taker_fee,omitempty"`       // taker fee, null if taker_side is not equal to side
	CreatedAt     string `json:"created_at,omitempty"`
}

// ReadTrades Trades of a asset pair, Only returns 50 latest trades
func ReadTrades(assetPairName string) ([]*Trade, error) {

	path := fmt.Sprintf("/asset_pairs/%v/trades", strings.ToUpper(assetPairName))

	resp, err := HTTPRequest(context.Background()).Get(path)

	if err != nil {
		return nil, err
	}

	trades := []*Trade{}

	if err := UnmarshalResponse(resp, &trades); err != nil {
		return nil, err
	}

	return trades, nil
}

// ReadUserTrades Trades of user
// Deleted
func ReadUserTrades(assetPairName string) ([]*Trade, error) {

	resp, err := HTTPRequest(context.Background()).SetQueryParam("asset_pair_name", assetPairName).Get("/viewer/trades")

	if err != nil {
		return nil, err
	}

	userTrades := []*Trade{}

	if err := UnmarshalResponse(resp, &userTrades); err != nil {
		return nil, err
	}

	return userTrades, nil
}
