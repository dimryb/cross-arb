package scan

import (
	"context"
	"log"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// DEXPriceUseCase — входной юзкейс для периодической загрузки цен с DEX.
// Периодически получает объем-зависимые котировки от DEX-провайдеров.
// Адаптеры обязаны реализовывать интерфейс DEXAdapter и предоставлять объем-зависимый метод Quote.
// Каналом out владеет вызывающая сторона. Реализация обязана уважать ctx.
//
// Формат выхода — ExecutableQuote: bid/ask — эффективные цены в QUOTE за BASE для указанного baseAmount.
// Для объем-зависимых DEX BidQty/AskQty устанавливаются равными baseAmount (в BASE).
type DEXPriceUseCase struct{}

// NewDEXPriceUseCase конструктор.
func NewDEXPriceUseCase() *DEXPriceUseCase { return &DEXPriceUseCase{} }

// Stream раз в interval обходит все пары и провайдеры, вызывает Quote и публикует ExecutableQuote в out.
// Bid/Ask — цены QUOTE за 1 BASE, BidQty/AskQty = baseAmount, Timestamp = now, Exchange = provider.Name().
func (u *DEXPriceUseCase) Stream(
	ctx context.Context,
	providers []i.DEXAdapter,
	pairs []string,
	interval time.Duration,
	baseAmount float64,
	out chan<- entity.ExecutableQuote,
) error {
	if ctx == nil {
		return context.Canceled
	}
	if interval <= 0 {
		interval = time.Second
	}

	// helper для одного прохода
	// helper для одного прохода
	scanOnce := func(now time.Time) {
		log.Printf("[DEX] tick now=%s providers=%d pairs=%d baseAmount=%.8f",
			now.Format(time.RFC3339), len(providers), len(pairs), baseAmount)

		for _, p := range providers {
			for _, pair := range pairs {
				log.Printf("[DEX] quote start ex=%s pair=%s base=%.8f", p.Name(), pair, baseAmount)
				bid, ask, err := p.Quote(ctx, pair, baseAmount)
				if err != nil {
					log.Printf("[DEX] quote ERR ex=%s pair=%s: %v", p.Name(), pair, err)
					continue
				}
				q := entity.ExecutableQuote{
					Exchange:  p.Name(),
					Pair:      pair,
					Bid:       bid,
					Ask:       ask,
					BidQty:    baseAmount,
					AskQty:    baseAmount,
					Timestamp: now,
				}
				select {
				case <-ctx.Done():
					log.Printf("[DEX] publish canceled ex=%s pair=%s", p.Name(), pair)
					return
				case out <- q:
					log.Printf("[DEX] publish OK ex=%s pair=%s bid=%.8f ask=%.8f", p.Name(), pair, bid, ask)
				}
			}
		}
	}

	// Немедленный старт
	scanOnce(time.Now())

	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-t.C:
			scanOnce(now)
		}
	}
}
