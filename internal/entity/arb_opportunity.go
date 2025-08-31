package entity

import "time"

// ArbOpportunity описывает арбитражную возможность между биржами.
// Содержит места покупки/продажи, соответствующие цены, валовую/чистую прибыль
// и относительный спред. Значения NetPnl/SpreadPct предполагают учёт комиссий.
type ArbOpportunity struct {
	Pair       string
	BuyOn      string
	BuyPrice   float64
	SellOn     string
	SellPrice  float64
	GrossPnl   float64 // Разница в цене без учёта комиссий
	NetPnl     float64 // Прибыль после вычета комиссий
	SpreadPct  float64 // NetPnl / BuyPrice, в процентах
	DetectedAt time.Time
}
