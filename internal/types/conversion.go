package types

import (
	"fmt"

	"github.com/dimryb/cross-arb/proto"
)

// ToBookTicker конвертирует proto.TickerData в BookTicker.
func ToBookTicker(pb *proto.TickerData) BookTicker {
	if pb == nil {
		return BookTicker{}
	}
	return BookTicker{
		Symbol:   pb.GetSymbol(),
		BidPrice: fmt.Sprintf("%f", pb.GetBidPrice()),
		BidQty:   fmt.Sprintf("%f", pb.GetBidQty()),
		AskPrice: fmt.Sprintf("%f", pb.GetAskPrice()),
		AskQty:   fmt.Sprintf("%f", pb.GetAskQty()),
	}
}
