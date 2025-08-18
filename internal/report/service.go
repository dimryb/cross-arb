package report

import (
	"fmt"
	"sync"
	"time"

	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/types"
)

// TickerBySymbolAndExchange — агрегированное хранилище тикеров: symbol → exchange → TickerData.
type TickerBySymbolAndExchange map[string]map[string]types.TickerData

type Service struct {
	log      i.Logger
	store    i.TickerStore
	sub      i.TickerSubscriber
	lastData TickerBySymbolAndExchange
	mu       sync.Mutex
}

func NewReportService(log i.Logger, store i.TickerStore) *Service {
	return &Service{
		log:      log,
		store:    store,
		lastData: make(TickerBySymbolAndExchange),
	}
}

func (r *Service) Start() {
	r.sub = r.store.AddSubscriber()

	go r.run()
}

func (r *Service) run() {
	eventCh := make(chan types.TickerEvent, 10)

	go func() {
		for {
			event, ok := r.sub.Recv()
			if !ok {
				close(eventCh)
				return
			}
			select {
			case eventCh <- event:
			default:
			}
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				r.log.Warnf("ReportService: event channel closed")
				return
			}
			r.handleEvent(event)

		case <-ticker.C:
			// r.printReport()
		}
	}
}

func (r *Service) handleEvent(event types.TickerEvent) {
	ticker := event.Ticker
	symbol := ticker.Symbol
	exchange := ticker.Exchange

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.lastData[symbol]; !ok {
		r.lastData[symbol] = make(map[string]types.TickerData)
	}
	r.lastData[symbol][exchange] = ticker
}

func (r *Service) printReport() { //nolint: unused
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.lastData) == 0 {
		return
	}

	timestamp := time.Now().Format("15:04:05.000")
	fmt.Printf("\n=== Обновление цен (%s) ===\n", timestamp)

	for symbol, exchanges := range r.lastData {
		for exchange, ticker := range exchanges {
			fmt.Printf(
				"  [%s @ %s] -> покупка: %.6f (%.4f) | продажа: %.6f (%.4f)\n",
				symbol, exchange,
				ticker.BidPrice, ticker.BidQty,
				ticker.AskPrice, ticker.AskQty,
			)
		}
	}
}

// Stop — корректно завершает подписку.
func (r *Service) Stop() {
	if r.sub != nil {
		r.sub.Close()
	}
}
