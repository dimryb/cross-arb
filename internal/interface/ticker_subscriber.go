package interfaces

import (
	"github.com/dimryb/cross-arb/internal/types"
)

// TickerSubscriber — абстракция канала получения событий.
type TickerSubscriber interface {
	Recv() (types.TickerEvent, bool) // (event, ok)
	Done() <-chan struct{}           // Контекст завершения
	Close()                          // Закрыть подписку
}
