package scan

import (
	"context"
	"errors"
	"sync"
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

// DEXPriceUseCaseImpl — базовая реализация Stream c ограничением конкуренции и таймаутом вызова адаптера.
// Потокобезопасна; уважает ctx; не владеет каналом out.
type DEXPriceUseCaseImpl struct {
	logger         i.Logger
	maxConcurrency int
	requestTimeout time.Duration
}

// NewDEXPriceUseCaseImpl создаёт реализацию DEX потока цен.
// maxConcurrency <= 0 заменяется на 8. requestTimeout == 0 означает использовать только родительский ctx.
func NewDEXPriceUseCaseImpl(l i.Logger, maxConcurrency int, requestTimeout time.Duration) *DEXPriceUseCaseImpl {
	if maxConcurrency <= 0 {
		maxConcurrency = 8
	}
	return &DEXPriceUseCaseImpl{
		logger:         l,
		maxConcurrency: maxConcurrency,
		requestTimeout: requestTimeout,
	}
}

// Stream периодически публикует котировки от каждого DEX-провайдера по каждой паре.
// На ошибках адаптера пишет в лог и продолжает поток; некорректные (<=0) цены отбрасываются.
func (u *DEXPriceUseCaseImpl) Stream(
	ctx context.Context,
	providers []i.DEXAdapter,
	pairs []string,
	interval time.Duration,
	baseAmount float64,
	out chan<- entity.ExecutableQuote,
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
	if baseAmount <= 0 {
		return errors.New("baseAmount не может быть отрицательным")
	}

	// Первый прогон сразу
	u.scanOnce(ctx, providers, pairs, baseAmount, out)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			u.scanOnce(ctx, providers, pairs, baseAmount, out)
			// Сбрасываем накопившиеся тики, если цикл занял больше interval
		drain:
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

func (u *DEXPriceUseCaseImpl) scanOnce(
	ctx context.Context,
	providers []i.DEXAdapter,
	pairs []string,
	baseAmount float64,
	out chan<- entity.ExecutableQuote,
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

				bid, ask, err := provider.Quote(callCtx, pair, baseAmount)
				if err != nil {
					u.logger.Warnf("DEX quote failed: ex=%s pair=%s err=%v", provider.Name(), pair, err)
					return
				}
				if bid <= 0 || ask <= 0 {
					u.logger.Warnf("DEX returned non-positive price: ex=%s pair=%s bid=%.8f ask=%.8f", provider.Name(), pair, bid, ask)
					return
				}

				q := entity.ExecutableQuote{
					Exchange:  provider.Name(),
					Pair:      pair,
					Bid:       bid,
					Ask:       ask,
					BidQty:    baseAmount, // для DEX — объём совпадает с запрошенным baseAmount
					AskQty:    baseAmount,
					Timestamp: time.Now(),
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
