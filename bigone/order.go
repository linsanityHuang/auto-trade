package bigone

import (
	"auto-trade/global"
	"context"
	"fmt"
)

// Side BID 买入, ASK 卖出
// Type LIMIT/MARKET/STOP_LIMIT/STOP_MARKET 限价, 市价, 限价触发, 市价触发

// Order Order
type Order struct {
	ID                int64  `json:"id,omitempty"`                  // id of order
	AssetPairName     string `json:"asset_pair_name,omitempty"`     // name of asset pair
	Price             string `json:"price,omitempty"`               // order price
	Amount            string `json:"amount,omitempty"`              // order amount
	FilledAmount      string `json:"filled_amount,omitempty"`       // already filled amount
	AvgDealPrice      string `json:"avg_deal_price,omitempty"`      // average price of the deal
	Side              string `json:"side,omitempty"`                // order side, one of ASK/BID
	State             string `json:"state,omitempty"`               // order status, one of FILLED/PENDING/CANCELLED/FIRED/REJECTED
	CreatedAt         string `json:"created_at,omitempty"`          // created time
	UpdatedAt         string `json:"updated_at,omitempty"`          // updated time
	Type              string `json:"type,omitempty"`                // order type, one of LIMIT/MARKET/STOP_LIMIT/STOP_MARKET
	StopPrice         string `json:"stop_price,omitempty"`          // order stop_price, only for stop order
	Operator          string `json:"operator,omitempty"`            // operator, one of GTE/LTE
	ImmediateOrCancel bool   `json:"immediate_or_cancel,omitempty"` // only use in limit type
	PostOnly          bool   `json:"post_only,omitempty"`           // only use in limit type
}

// ReadOrders Get user orders in a asset pair
func ReadOrders(assetPairName string) ([]*Order, error) {
	token, err := SignAuthenticationToken(global.BigOneSetting.APIKEY, global.BigOneSetting.APISECRET)
	if err != nil {
		return nil, err
	}

	resp, err := HTTPRequest(context.Background()).
		SetAuthToken(token).
		SetQueryParam("asset_pair_name", assetPairName).
		Get("/viewer/orders")

	if err != nil {
		return nil, err
	}

	orders := []*Order{}

	if err := UnmarshalResponse(resp, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

// ReadOrder Get one order
func ReadOrder(orderID int64) (*Order, error) {

	token, err := SignAuthenticationToken(global.BigOneSetting.APIKEY, global.BigOneSetting.APISECRET)
	if err != nil {
		return nil, err
	}

	resp, err := HTTPRequest(context.Background()).
		SetAuthToken(token).
		Get(fmt.Sprintf("/viewer/orders/%v", orderID))

	if err != nil {
		return nil, err
	}

	order := &Order{}

	if err := UnmarshalResponse(resp, order); err != nil {
		return nil, err
	}

	return order, nil
}

// CreateOrder Create Order
func CreateOrder(assetPairName, side, price, amount, orderType string) (*Order, error) {
	token, err := SignAuthenticationToken(global.BigOneSetting.APIKEY, global.BigOneSetting.APISECRET)
	if err != nil {
		return nil, err
	}

	order := Order{
		AssetPairName: assetPairName,
		Side:          side,
		Price:         price,
		Amount:        amount,
		Type:          orderType,
	}

	resp, err := HTTPRequest(context.Background()).
		SetAuthToken(token).
		SetBody(order).
		Post("/viewer/orders")

	if err != nil {
		return nil, err
	}

	o := &Order{}

	if err = UnmarshalResponse(resp, o); err != nil {
		return nil, err
	}

	return o, nil
}

// CreateMultiOrder Multiple Create Orders
func CreateMultiOrder(orders []*Order) ([]*Order, error) {
	// Notice: If an error occurs in it, all order creations fail.
	// The quantity limit of multiple create orders: 10

	token, err := SignAuthenticationToken(global.BigOneSetting.APIKEY, global.BigOneSetting.APISECRET)
	if err != nil {
		return nil, err
	}

	resp, err := HTTPRequest(context.Background()).
		SetAuthToken(token).
		SetBody(orders).
		Post("/viewer/orders/multi")

	if err != nil {
		return nil, err
	}

	ors := []*Order{}

	if err = UnmarshalResponse(resp, &ors); err != nil {
		return nil, err
	}

	return ors, nil
}

// CancelOrder cancel order
func CancelOrder(orderID int64) (*Order, error) {
	resp, err := HTTPRequest(context.Background()).Post(fmt.Sprintf("/viewer/orders/%v/cancel", orderID))
	if err != nil {
		return nil, err
	}

	order := &Order{}

	if err := UnmarshalResponse(resp, order); err != nil {
		return nil, err
	}

	return order, nil
}

// CancelAllOrders Cancel All Orders
func CancelAllOrders() {

}
