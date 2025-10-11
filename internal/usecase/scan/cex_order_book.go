package scan

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// CEXOrderBookUseCase периодически загружает стаканы с CEX.
// Адаптеры обязаны реализовывать интерфейс CEXAdapter и предоставлять OrderBookDepth.
// Каналом out владеет вызывающая сторона. Реализация обязана уважать ctx.
//
// Формат выхода — OrderBookResult: символ (пара), данные стакана и возможная ошибка запроса.
// Реализация может выбирать политику лимита глубины (limit) — через параметры конструктора.
type CEXOrderBookUseCase interface {
	Stream(
		ctx context.Context,
		providers []i.CEXAdapter,
		pairs []string,
		interval time.Duration,
		out chan<- entity.OrderBookResult,
	) error
}

// CEXOrderBookUseCaseImpl — заглушка, полезна на этапе интеграции.
type CEXOrderBookUseCaseImpl struct {
	logger         i.Logger
	maxConcurrency int
	requestTimeout time.Duration
	depthLimit     int
}

func NewCEXOrderBookUseCaseImpl(
	l i.Logger,
	maxConcurrency int,
	requestTimeout time.Duration,
	depthLimit int,
) *CEXOrderBookUseCaseImpl {
	return &CEXOrderBookUseCaseImpl{
		logger:         l,
		maxConcurrency: maxConcurrency,
		requestTimeout: requestTimeout,
		depthLimit:     depthLimit,
	}
}

func (u *CEXOrderBookUseCaseImpl) Stream(
	ctx context.Context,
	providers []i.CEXAdapter,
	pairs []string,
	interval time.Duration,
	out chan<- entity.OrderBookResult,
) error {
	if out == nil {
		return errors.New("канал out не может быть nil")
	}
	if len(providers) == 0 {
		return errors.New("отсутствуют DEX providers")
	}
	if len(pairs) == 0 {
		return errors.New("отсутствуют pairs")
	}
	if interval <= 0 {
		return errors.New("interval должен быть положительным")
	}

	u.scanOnce(ctx, providers, pairs, out)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			u.scanOnce(ctx, providers, pairs, out)
		drain:
			// Сбрасываем накопившиеся тики, если цикл занял больше interval
			for {
				select {
				case <-ticker.C:
					continue drain
				default:
					break drain
				}
			}
		}
	}
}

func (u *CEXOrderBookUseCaseImpl) scanOnce(
	ctx context.Context,
	providers []i.CEXAdapter,
	pairs []string,
	out chan<- entity.OrderBookResult,
) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, u.maxConcurrency)

	for _, provider := range providers {
		for _, pair := range pairs {
			wg.Add(1)
			go func() {
				defer wg.Done()
				// Ограничение конкуренции
				sem <- struct{}{}
				defer func() { <-sem }()

				callCtx := ctx
				var cancel context.CancelFunc
				if u.requestTimeout > 0 {
					callCtx, cancel = context.WithTimeout(ctx, u.requestTimeout)
					defer cancel()
				}

				orderBook, err := provider.OrderBookDepth(callCtx, pair, u.depthLimit)
				if err != nil {
					u.logger.Warnf("CEX OrderBookDepth failed: ex=%s pair=%s err=%v", provider.Name(), pair, err)
					return
				}
				if len(orderBook.Bids) == 0 || len(orderBook.Asks) == 0 {
					u.logger.Warnf("CEX returned empty order book: ex=%s pair=%s", provider.Name(), pair)
					return
				}

				q := entity.OrderBookResult{
					Symbol: pair,
					Data:   orderBook,
					Error:  nil,
				}

				select {
				case out <- q:
				case <-ctx.Done():
					return
				}
			}()
		}
	}
	wg.Wait()
}

// NoopCEXOrderBookUseCase — заглушка, полезна на этапе интеграции.
type NoopCEXOrderBookUseCase struct{}

func NewNoopCEXOrderBookUseCase() *NoopCEXOrderBookUseCase { return &NoopCEXOrderBookUseCase{} }

func (n *NoopCEXOrderBookUseCase) Stream(
	ctx context.Context,
	_ []i.CEXAdapter,
	_ []string,
	_ time.Duration,
	_ chan<- entity.OrderBookResult,
) error {
	<-ctx.Done()
	return ctx.Err()
}
