package request

type Buy struct {
	Token    string  `validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
}

type Sell struct {
	Token    string  `validate:"required"`
	Amount   float64 `json:"amount" validate:"required"`
	Currency string  `json:"currency" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
}

type Webhook struct {
	Password  string  `json:"password" validate:"required"`
	Operation string  `json:"operation" validate:"required"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
	Price     float64 `json:"price"`
}
