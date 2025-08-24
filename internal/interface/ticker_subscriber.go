package interfaces

import (
	"github.com/dimryb/cross-arb/internal/entity"
)

//go:generate mockgen -source=ticker_subscriber.go -package=mocks -destination=../../mocks/mock_ticker_subscriber.go

// TickerSubscriber — абстракция канала получения событий.
type TickerSubscriber interface {
	Recv() (entity.TickerEvent, bool) // (event, ok)
	Done() <-chan struct{}            // Контекст завершения
	Close()                           // Закрыть подписку
}
