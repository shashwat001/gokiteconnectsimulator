package kiteconnectsimulator

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"main/kiteconnectsimulator/models"

	"github.com/google/go-querystring/query"
)

// Orders is a list of orders.
type Orders []models.Order

// OrderParams represents parameters for placing an order.
type OrderParams struct {
	Exchange        string `url:"exchange,omitempty"`
	Tradingsymbol   string `url:"tradingsymbol,omitempty"`
	Validity        string `url:"validity,omitempty"`
	ValidityTTL     int    `url:"validity_ttl,omitempty"`
	Product         string `url:"product,omitempty"`
	OrderType       string `url:"order_type,omitempty"`
	TransactionType string `url:"transaction_type,omitempty"`

	Quantity          int     `url:"quantity,omitempty"`
	DisclosedQuantity int     `url:"disclosed_quantity,omitempty"`
	Price             float64 `url:"price,omitempty"`
	TriggerPrice      float64 `url:"trigger_price,omitempty"`

	Squareoff        float64 `url:"squareoff,omitempty"`
	Stoploss         float64 `url:"stoploss,omitempty"`
	TrailingStoploss float64 `url:"trailing_stoploss,omitempty"`

	IcebergLegs int `url:"iceberg_legs,omitempty"`
	IcebergQty  int `url:"iceberg_quantity,omitempty"`

	Tag string `json:"tag" url:"tag,omitempty"`
}

// OrderResponse represents the order place success response.
type OrderResponse struct {
	OrderID string `json:"order_id"`
}

// Trade represents an individual trade response.
type Trade struct {
	AveragePrice      float64     `json:"average_price"`
	Quantity          float64     `json:"quantity"`
	TradeID           string      `json:"trade_id"`
	Product           string      `json:"product"`
	FillTimestamp     models.Time `json:"fill_timestamp"`
	ExchangeTimestamp models.Time `json:"exchange_timestamp"`
	ExchangeOrderID   string      `json:"exchange_order_id"`
	OrderID           string      `json:"order_id"`
	TransactionType   string      `json:"transaction_type"`
	TradingSymbol     string      `json:"tradingsymbol"`
	Exchange          string      `json:"exchange"`
	InstrumentToken   uint32      `json:"instrument_token"`
}

// Trades is a list of trades.
type Trades []Trade

// GetOrders gets list of orders.
func (c *Client) GetOrders() (Orders, error) {
	var orders Orders
	err := c.doEnvelope(http.MethodGet, URIGetOrders, nil, nil, &orders)
	return orders, err
}

// GetTrades gets list of trades.
func (c *Client) GetTrades() (Trades, error) {
	var trades Trades
	err := c.doEnvelope(http.MethodGet, URIGetTrades, nil, nil, &trades)
	return trades, err
}

// GetOrderHistory gets history of an individual order.
func (c *Client) GetOrderHistory(OrderID string) ([]models.Order, error) {
	var orderHistory []models.Order
	err := c.doEnvelope(http.MethodGet, fmt.Sprintf(URIGetOrderHistory, OrderID), nil, nil, &orderHistory)
	return orderHistory, err
}

// GetOrderTrades gets list of trades executed for a particular order.
func (c *Client) GetOrderTrades(OrderID string) ([]Trade, error) {
	var orderTrades []Trade
	err := c.doEnvelope(http.MethodGet, fmt.Sprintf(URIGetOrderTrades, OrderID), nil, nil, &orderTrades)
	return orderTrades, err
}

// PlaceOrder places an order.
func (c *Client) PlaceOrder(variety string, orderParams OrderParams) (OrderResponse, error) {
	log.Println("Adding order")
	var (
		orderResponse OrderResponse
		err           error
	)

	if _, err = query.Values(orderParams); err != nil {
		return orderResponse, NewError(InputError, fmt.Sprintf("Error decoding order params: %v", err), nil)
	}

	instrument := fmt.Sprintf("%s:%s", orderParams.Exchange, orderParams.Tradingsymbol)
	quoteLTP, err := c.GetLTP(instrument)
	if err != nil {
		log.Fatal("Error in getting quote for instrument: ", err)
	}

	order := &DbOrder{Order: models.Order{
		Exchange:        orderParams.Exchange,
		TradingSymbol:   orderParams.Tradingsymbol,
		Quantity:        float64(orderParams.Quantity),
		Price:           quoteLTP[instrument].LastPrice,
		InstrumentToken: uint32(quoteLTP[instrument].InstrumentToken),
		TransactionType: orderParams.TransactionType,
		OrderType:       orderParams.OrderType,
		Status:          OrderStatusOpen,
	}}

	if orderParams.OrderType == OrderTypeLimit {
		order.Status = OrderStatusOpen
	} else {
		order.Status = OrderStatusComplete
	}

	_, err = db.NewInsert().Model(order).Exec(context.Background())

	if err != nil {
		panic(err)
	}

	if order.OrderType == OrderTypeMarket {
		complete_order_and_update_holding(order.ID)
	} else if order.OrderType == OrderTypeLimit {
		if order.TransactionType == TransactionTypeBuy {
			c.Om.AddBuy(instrument, order.InstrumentToken, int64(order.Quantity), order.Price)
		} else if order.TransactionType == TransactionTypeSell {
			c.Om.AddSell(instrument, order.InstrumentToken, int64(order.Quantity), order.Price)
		} else {
			panic("Unknown transactiontype")
		}
	}

	return OrderResponse{strconv.FormatInt(order.ID, 10)}, err
}

// ModifyOrder modifies an order.
func (c *Client) ModifyOrder(variety string, orderID string, orderParams OrderParams) (OrderResponse, error) {
	var (
		orderResponse OrderResponse
		params        url.Values
		err           error
	)

	if params, err = query.Values(orderParams); err != nil {
		return orderResponse, NewError(InputError, fmt.Sprintf("Error decoding order params: %v", err), nil)
	}

	err = c.doEnvelope(http.MethodPut, fmt.Sprintf(URIModifyOrder, variety, orderID), params, nil, &orderResponse)
	return orderResponse, err
}

// CancelOrder cancels/exits an order.
func (c *Client) CancelOrder(variety string, orderID string, parentOrderID *string) (OrderResponse, error) {
	var (
		orderResponse OrderResponse
		params        url.Values
	)

	if parentOrderID != nil {
		// initialize the params map first
		params := url.Values{}
		params.Add("parent_order_id", *parentOrderID)
	}

	err := c.doEnvelope(http.MethodDelete, fmt.Sprintf(URICancelOrder, variety, orderID), params, nil, &orderResponse)
	return orderResponse, err
}

// ExitOrder is an alias for CancelOrder which is used to cancel/exit an order.
func (c *Client) ExitOrder(variety string, orderID string, parentOrderID *string) (OrderResponse, error) {
	return c.CancelOrder(variety, orderID, parentOrderID)
}
