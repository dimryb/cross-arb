package scan

import (
	"context"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// DEXPriceUseCase — входной юзкейс для периодической загрузки цен с DEX.
// Адаптеры обязаны реализовывать интерфейс DEXAdapter и предоставлять объем-зависимый метод Quote.
// Каналом out владеет вызывающая сторона. Реализация обязана уважать ctx.
//
// Формат выхода — ExecutableQuote: bid/ask — эффективные цены в QUOTE за BASE для указанного baseAmount.
// Для объем-зависимых DEX BidQty/AskQty устанавливаются равными baseAmount (в BASE).
type DEXPriceUseCase interface {
	Stream(
		ctx context.Context,
		providers []i.DEXAdapter,
		pairs []string,
		interval time.Duration,
		baseAmount float64,
		out chan<- entity.ExecutableQuote,
	) error
}

// NoopDEXPriceUseCase — заглушка, полезна для стадий интеграции.
type NoopDEXPriceUseCase struct{}

func NewNoopDEXPriceUseCase() *NoopDEXPriceUseCase { return &NoopDEXPriceUseCase{} }

func (n *NoopDEXPriceUseCase) Stream(
	ctx context.Context,
	_ []i.DEXAdapter,
	_ []string,
	_ time.Duration,
	_ float64,
	_ chan<- entity.ExecutableQuote,
) error {
	<-ctx.Done()
	return ctx.Err()
}
