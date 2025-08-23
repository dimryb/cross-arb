package interfaces

import "github.com/dimryb/cross-arb/internal/entity"

//go:generate mockgen -source=tiсker_store.go -package=mocks -destination=../../mocks/mock_ticker_store.go
type TickerStore interface {
	Set(t entity.TickerData)
	GetAll() []entity.TickerData
	AddSubscriber() TickerSubscriber
}
