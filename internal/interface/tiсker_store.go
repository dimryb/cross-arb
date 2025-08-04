package interfaces

import "github.com/dimryb/cross-arb/internal/types"

//go:generate mockgen -source=tiсker_store.go -package=mocks -destination=../../mocks/mock_ticker_store.go
type TickerStore interface {
	Set(t types.TickerData)
	GetAll() []types.TickerData
	AddSubscriber() TickerSubscriber
}
