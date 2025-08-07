package scanner

import "time"

// PricePoint содержит котировки, полученные от биржи.
type PricePoint struct {
	Exchange  string    // Имя биржи (adapter.Name()).
	Pair      string    // Например, "BTC/USDT".
	Bid       float64   // Лучшая цена покупки.
	Ask       float64   // Лучшая цена продажи.
	Timestamp time.Time // Время получения данных сканером.
}

// Opportunity описывает найденную арбитражную ситуацию.
type Opportunity struct {
	Pair       string
	BuyOn      string
	BuyPrice   float64
	SellOn     string
	SellPrice  float64
	GrossPnl   float64 // Разница цен без учёта комиссий.
	NetPnl     float64 // Прибыль после вычета комиссий.
	SpreadPct  float64 // NetPnl / BuyPrice, процентов.
	DetectedAt time.Time
}
