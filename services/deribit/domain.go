package deribit

type Authentication struct {
	Result Auth `json:"result"`
}

type Auth struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Price struct {
	Currency []Currency `json:"result"`
}

type Currency struct {
	QuoteCurrency          string  `json:"quote_currency"`
	InstrumentName         string  `json:"instrument_name"`
	BaseCurrency           string  `json:"base_currency"`
	MarkPrice              float64 `json:"mark_price"`
	EstimatedDeliveryPrice float64 `json:"estimated_delivery_price"`
	PriceChange            float64 `json:"price_change"`
}

type Buy struct {
	ID     int64  `json:"id"`
	Result Result `json:"result"`
}

type Sell struct {
	Result Result `json:"result"`
}

type Result struct {
	Trades []Trade `json:"trades"`
	Order  Order   `json:"order"`
}

type Trade struct {
	TradeID        string  `json:"trade_id"`
	InstrumentName string  `json:"instrument_name"`
	Price          float64 `json:"price"`
	Amount         float64 `json:"amount"`
}

type Order struct {
	ProfitLoss     float64 `json:"profit_loss"`
	InstrumentName string  `json:"instrument_name"`
	Price          float64 `json:"price"`
	Amount         float64 `json:"amount"`
}
