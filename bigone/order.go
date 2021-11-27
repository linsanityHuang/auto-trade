package bigone

import (
	"auto-trade/global"
	"context"
	"fmt"
	"log"
	"time"
)

// Side
type Side int64

const (
	BID Side = iota // 买入
	ASK             // 卖出
)

func (s Side) String() string {
	switch s {
	case BID:
		return "BID"
	case ASK:
		return "ASK"
	}
	return "UNKNOWN"
}

type State int64

const (
	UNKNOWN State = iota
	FILLED
	PENDING
	CANCELLED
	FIRED
	REJECTED
)

func (s State) String() string {
	switch s {
	case FILLED:
		return "FILLED"
	case PENDING:
		return "PENDING"
	case CANCELLED:
		return "CANCELLED"
	case FIRED:
		return "FIRED"
	case REJECTED:
		return "REJECTED"
	}
	return "UNKNOWN"
}

type OrderType int64

const (
	LIMIT       OrderType = iota // 限价
	MARKET                       // 市价
	STOP_LIMIT                   // 限价触发
	STOP_MARKET                  // 市价触发
)

func (o OrderType) String() string {
	switch o {
	case LIMIT:
		return "LIMIT"
	case MARKET:
		return "MARKET"
	case STOP_LIMIT:
		return "STOP_LIMIT"
	case STOP_MARKET:
		return "STOP_MARKET"
	}
	return "UNKNOWN"
}

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
func ReadOrders(assetPairName string, state State) ([]*Order, error) {
	resp, err := HTTPRequest(context.Background()).
		SetQueryParam("asset_pair_name", assetPairName).
		SetQueryParam("state", state.String()).
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
func ReadOrder(id int64) (*Order, error) {
	resp, err := HTTPRequest(context.Background()).Get(fmt.Sprintf("/viewer/orders/%v", id))

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
// Price's scale should less than AssetPair's quote scale 价格的规模应小于资产对的报价规模
// Amount's scale should less than AssetPair's base scale 金额的规模应小于资产对的基本规模
// Price * Amount should larger than AssetPair's min_quote_value 价格*金额应大于资产对的最小报价值
func CreateOrder(assetPairName, side, price, amount, orderType string) (*Order, error) {

	order := Order{
		AssetPairName: assetPairName,
		Side:          side,
		Price:         price,
		Amount:        amount,
		Type:          orderType,
	}

	resp, err := HTTPRequest(context.Background()).SetBody(order).Post("/viewer/orders")

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
// Notice: If an error occurs in it, all order creations fail.
// The quantity limit of multiple create orders: 10
func CreateMultiOrder(orders []*Order) ([]*Order, error) {
	resp, err := HTTPRequest(context.Background()).
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
func CancelOrder(id int64) (*Order, error) {

	resp, err := HTTPRequest(context.Background()).
		Post(fmt.Sprintf("/viewer/orders/%v/cancel", id))

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
// TODO
func CancelAllOrders(assetPairName string) ([]int64, error) {
	resp, err := HTTPRequest(context.Background()).Post("/viewer/orders/cancel")

	if err != nil {
		return nil, err
	}

	var ids []int64

	if err := UnmarshalResponse(resp, ids); err != nil {
		return nil, err
	}

	return ids, nil
}

// WaitBidOrder 等待买入订单成交, 超时则取消订单
func (o *Order) WaitBidOrder(timeout time.Duration) bool {
	if o.Side != BID.String() {
		log.Fatalf("订单状态异常, 此订单的Side应为BID.")
	}

	start := time.Now()

	for time.Now().Sub(start).Minutes() <= timeout.Minutes() {
		order, err := ReadOrder(o.ID)
		if err != nil {
			log.Printf("read order err in wait order: %v\n", err)
			continue
		}

		if order.State != FILLED.String() {
			continue
		}

		return true
	}

	// 取消挂单
	_, err := CancelOrder(o.ID)
	if err != nil {
		log.Printf("cancel order err in wait order: %v\n", err)
	} else {
		log.Println("cancel bid order success.")
	}

	return false
}

// WaitBidOrderForever 永远等待买入订单成交
func (o *Order) WaitBidOrderForever() bool {
	if o.Side != BID.String() {
		log.Fatalf("订单状态异常, 此订单的Side应为BID.")
	}

	for {
		order, err := ReadOrder(o.ID)
		if err != nil {
			log.Printf("read order err in wait order: %v\n", err)
			continue
		}

		if order.State != FILLED.String() {
			continue
		}

		return true
	}
}

// WaitAskOrderForever 永远等待卖出订单成交
func (o *Order) WaitAskOrderForever() bool {
	if o.Side != ASK.String() {
		log.Fatalf("订单状态异常, 此订单的Side应为ASK.")
	}

	for {

		time.Sleep(global.Timeout)

		order, err := ReadOrder(o.ID)
		if err != nil {
			log.Printf("read order err in wait order: %v\n", err)
			continue
		}

		if order.State != FILLED.String() {
			continue
		}

		return true
	}
}

func ShowOrdes() {
	existOrders, err := ReadOrders(BTCUSDT, PENDING)
	if err != nil {
		fmt.Printf("read exist order failed: %v\n", err)
		return
	}

	for _, order := range existOrders {
		fmt.Printf("当前订单, ID: %d, AssetPairName: %s, Side: %s, State: %s, Price: %s, Amount: %s\n",
			order.ID, order.AssetPairName, order.Side, order.State, order.Price, order.Amount)
	}
}
