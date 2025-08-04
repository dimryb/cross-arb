package storage

import (
	"context"
	"sync"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/types"
)

// subscriber — приватная реализация i.TickerSubscriber.
type subscriber struct {
	ctx     context.Context
	eventCh <-chan types.TickerEvent
	cancel  context.CancelFunc
	once    sync.Once
}

// newSubscriber создаёт новую подписку.
func newSubscriber(
	ctx context.Context,
	eventCh <-chan types.TickerEvent,
	cancel context.CancelFunc,
) i.TickerSubscriber {
	return &subscriber{
		ctx:     ctx,
		eventCh: eventCh,
		cancel:  cancel,
	}
}

// Recv — возвращает types.TickerEvent, как в интерфейсе.
func (s *subscriber) Recv() (types.TickerEvent, bool) {
	select {
	case event, ok := <-s.eventCh:
		return event, ok
	case <-s.ctx.Done():
		var zero types.TickerEvent
		return zero, false
	}
}

func (s *subscriber) Done() <-chan struct{} {
	return s.ctx.Done()
}

func (s *subscriber) Close() {
	s.once.Do(s.cancel)
}
