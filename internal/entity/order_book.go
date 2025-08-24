package entity

type Order struct {
	Price    float64
	Quantity float64
}

type OrderBookResult struct {
	Symbol string
	Data   OrderBook
	Error  error
}

type OrderBook struct {
	Bids []Order `json:"bids"`
	Asks []Order `json:"asks"`
}
