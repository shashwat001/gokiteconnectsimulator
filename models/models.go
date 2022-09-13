package models

// Order represents a individual order response.
type Order struct {
	AccountID string `json:"account_id"`
	PlacedBy  string `json:"placed_by"`

	OrderID          string                 `json:"order_id"`
	ExchangeOrderID  string                 `json:"exchange_order_id"`
	ParentOrderID    string                 `json:"parent_order_id"`
	Status           string                 `json:"status"`
	StatusMessage    string                 `json:"status_message"`
	StatusMessageRaw string                 `json:"status_message_raw"`
	Variety          string                 `json:"variety"`
	Meta             map[string]interface{} `json:"meta"`

	Exchange        string `json:"exchange"`
	TradingSymbol   string `json:"tradingsymbol"`
	InstrumentToken uint32 `json:"instrument_token"`

	OrderType         string  `json:"order_type"`
	TransactionType   string  `json:"transaction_type"`
	Validity          string  `json:"validity"`
	ValidityTTL       int     `json:"validity_ttl"`
	Product           string  `json:"product"`
	Quantity          float64 `json:"quantity"`
	DisclosedQuantity float64 `json:"disclosed_quantity"`
	Price             float64 `json:"price"`
	TriggerPrice      float64 `json:"trigger_price"`

	AveragePrice      float64 `json:"average_price"`
	FilledQuantity    float64 `json:"filled_quantity"`
	PendingQuantity   float64 `json:"pending_quantity"`
	CancelledQuantity float64 `json:"cancelled_quantity"`

	Tag  string   `json:"tag"`
	Tags []string `json:"tags"`
}

// Holding is an individual holdings response.
type Holding struct {
	Tradingsymbol   string `json:"tradingsymbol"`
	Exchange        string `json:"exchange"`
	InstrumentToken uint32 `json:"instrument_token"`
	ISIN            string `json:"isin"`
	Product         string `json:"product"`

	Price              float64 `json:"price"`
	UsedQuantity       int     `json:"used_quantity"`
	Quantity           int     `json:"quantity"`
	T1Quantity         int     `json:"t1_quantity"`
	RealisedQuantity   int     `json:"realised_quantity"`
	AuthorisedQuantity int     `json:"authorised_quantity"`
	OpeningQuantity    int     `json:"opening_quantity"`
	CollateralQuantity int     `json:"collateral_quantity"`
	CollateralType     string  `json:"collateral_type"`

	Discrepancy         bool    `json:"discrepancy"`
	AveragePrice        float64 `json:"average_price"`
	LastPrice           float64 `json:"last_price"`
	ClosePrice          float64 `json:"close_price"`
	PnL                 float64 `json:"pnl"`
	DayChange           float64 `json:"day_change"`
	DayChangePercentage float64 `json:"day_change_percentage"`
}
