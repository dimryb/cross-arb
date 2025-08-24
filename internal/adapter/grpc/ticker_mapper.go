package grpc

import (
	"fmt"

	"github.com/dimryb/cross-arb/internal/entity"
	"github.com/dimryb/cross-arb/proto"
)

// ToBookTicker конвертирует proto.TickerData в BookTicker.
func ToBookTicker(pb *proto.TickerData) entity.BookTicker {
	if pb == nil {
		return entity.BookTicker{}
	}
	return entity.BookTicker{
		Symbol:   pb.GetSymbol(),
		BidPrice: fmt.Sprintf("%f", pb.GetBidPrice()),
		BidQty:   fmt.Sprintf("%f", pb.GetBidQty()),
		AskPrice: fmt.Sprintf("%f", pb.GetAskPrice()),
		AskQty:   fmt.Sprintf("%f", pb.GetAskQty()),
	}
}
