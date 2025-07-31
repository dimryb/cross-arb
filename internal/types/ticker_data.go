package types

// TickerData — данные одного тикера.
type TickerData struct {
	Symbol   string  `json:"symbol" proto:"symbol"`
	Exchange string  `json:"exchange" proto:"exchange"`
	BidPrice float64 `json:"bidPrice" proto:"bid_price"`
	BidQty   float64 `json:"bidQty" proto:"bid_qty"`
	AskPrice float64 `json:"askPrice" proto:"ask_price"`
	AskQty   float64 `json:"askQty" proto:"ask_qty"`
}

// Equal Сравнение двух тикеров (для оптимизации уведомлений).
func (t TickerData) Equal(other TickerData) bool {
	return t.Symbol == other.Symbol &&
		t.Exchange == other.Exchange &&
		t.BidPrice == other.BidPrice &&
		t.BidQty == other.BidQty &&
		t.AskPrice == other.AskPrice &&
		t.AskQty == other.AskQty
}
