package scan

import (
	"context"
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
