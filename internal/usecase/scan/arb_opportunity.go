package scan

import (
	"context"
	"math"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// ArbOpportunityUseCase хранит последние котировки и ищет арбитражные возможности между биржами с учётом комиссий.
type ArbOpportunityUseCase struct {
	adapters map[string]i.EXAdapter // по имени биржи
}

// NewOpportunityUseCase принимает список адаптеров, чтобы уметь узнавать их комиссии.
func NewOpportunityUseCase(adapters []i.EXAdapter) *ArbOpportunityUseCase {
	m := make(map[string]i.EXAdapter, len(adapters))
	for _, a := range adapters {
		if a == nil {
			continue
		}
		m[a.Name()] = a
	}
	return &ArbOpportunityUseCase{adapters: m}
}

// Detect читает поток котировок и публикует арбитражные возможности при наличии положительного чистого спреда.
func (u *ArbOpportunityUseCase) Detect(
	ctx context.Context,
	in <-chan entity.ExecutableQuote,
	out chan<- entity.ArbOpportunity,
) error {
	if ctx == nil {
		return context.Canceled
	}
	// last[pair][exchange] = quote
	last := make(map[string]map[string]entity.ExecutableQuote)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case q, ok := <-in:
			if !ok {
				return nil // канал закрыт — завершаем без ошибки
			}

			// Инициализируем мапу для пары, если ещё не существует
			if _, ok := last[q.Pair]; !ok {
				last[q.Pair] = make(map[string]entity.ExecutableQuote)
			}
			last[q.Pair][q.Exchange] = q

			// Пытаемся найти арбитражную возможность для этой пары
			if opp := u.DetectOpportunityForPair(q.Pair, last[q.Pair]); opp != nil {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case out <- *opp:
				}
			}
		}
	}
}

// DetectOpportunityForPair анализирует котировки по паре и возвращает арбитражную возможность, если она есть.
// Возвращает nil, если арбитража нет.
func (u *ArbOpportunityUseCase) DetectOpportunityForPair(
	pair string,
	quotes map[string]entity.ExecutableQuote,
) *entity.ArbOpportunity {
	var (
		bestBuyEx  string
		bestBuy    = math.Inf(1)
		bestSellEx string
		bestSell   = math.Inf(-1)
	)

	// Находим лучшие цены покупки (ask) и продажи (bid)
	for ex, qq := range quotes {
		if qq.Ask > 0 && qq.Ask < bestBuy {
			bestBuy, bestBuyEx = qq.Ask, ex
		}
		if qq.Bid > bestSell {
			bestSell, bestSellEx = qq.Bid, ex
		}
	}

	// Проверяем, что нашли обе цены и на разных биржах
	if bestBuyEx == "" || bestSellEx == "" || bestBuyEx == bestSellEx {
		return nil
	}

	// Получаем комиссии
	buyTaker := u.GetTakerFee(bestBuyEx, pair)
	sellTaker := u.GetTakerFee(bestSellEx, pair)

	// Эффективные цены с учётом комиссий
	effBuy := bestBuy * (1 + buyTaker)
	effSell := bestSell * (1 - sellTaker)
	gross := bestSell - bestBuy
	net := effSell - effBuy

	// Арбитраж есть только если чистая прибыль положительна
	if net <= 0 {
		return nil
	}

	// Формируем объект арбитражной возможности
	return &entity.ArbOpportunity{
		Pair:       pair,
		BuyOn:      bestBuyEx,
		BuyPrice:   bestBuy,
		SellOn:     bestSellEx,
		SellPrice:  bestSell,
		GrossPnl:   gross,
		NetPnl:     net,
		SpreadPct:  (net / bestBuy) * 100,
		DetectedAt: time.Now(),
	}
}

// getTakerFee возвращает комиссию тейкера для указанной биржи и пары.
// Если адаптер не найден — возвращает 0.
func (u *ArbOpportunityUseCase) GetTakerFee(exchange, pair string) float64 {
	if a, ok := u.adapters[exchange]; ok {
		_, taker := a.TradingFee(pair)
		return taker
	}
	return 0.0
}
