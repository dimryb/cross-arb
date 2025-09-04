package scan

import (
	"context"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// CEXOrderBookUseCase периодически получает стаканы от CEX-провайдеров.
// Адаптеры обязаны реализовывать интерфейс CEXAdapter и предоставлять OrderBookDepth.
// Каналом out владеет вызывающая сторона. Реализация обязана уважать ctx.
//
// Формат выхода — OrderBookResult: символ (пара), данные стакана и возможная ошибка запроса.
// Реализация может выбирать политику лимита глубины (limit) — через параметры конструктора.
type CEXOrderBookUseCase struct {
	limit int
}

// NewCEXOrderBookUseCase с ограничением глубины стакана (если 0 — провайдер решит сам или возьмем 5).
func NewCEXOrderBookUseCase(limit int) *CEXOrderBookUseCase {
	return &CEXOrderBookUseCase{limit: limit}
}

// Stream раз в interval обходит все пары и провайдеры, публикуя OrderBookResult.
func (u *CEXOrderBookUseCase) Stream(
	ctx context.Context,
	providers []i.CEXAdapter,
	pairs []string,
	interval time.Duration,
	out chan<- entity.OrderBookResult,
) error {
	if ctx == nil {
		return context.Canceled
	}
	if interval <= 0 {
		interval = time.Second
	}
	limit := u.limit
	if limit < 0 {
		limit = 0
	}

	scanOnce := func() {
		for _, p := range providers {
			for _, pair := range pairs {
				book, err := p.OrderBookDepth(ctx, pair, limit)
				res := entity.OrderBookResult{Symbol: pair, Data: book, Error: err}
				select {
				case <-ctx.Done():
					return
				case out <- res:
				}
			}
		}
	}

	// Немедленный старт
	scanOnce()

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			scanOnce()
		}
	}
}
