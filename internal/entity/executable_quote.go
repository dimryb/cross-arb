package entity

import "time"

// ExecutableQuote представляет исполнимую котировку для пары.
// Значения Bid/Ask выражены в QUOTE за 1 BASE. Для DEX котировка рассчитывается
// под заданный объём, для CEX может агрегировать уровни стакана для покрытия объёма.
type ExecutableQuote struct {
	Exchange  string
	Pair      string
	Bid       float64
	Ask       float64
	BidQty    float64
	AskQty    float64
	Timestamp time.Time
}
