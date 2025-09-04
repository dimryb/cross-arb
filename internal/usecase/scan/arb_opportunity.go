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
			if _, ok := last[q.Pair]; !ok {
				last[q.Pair] = make(map[string]entity.ExecutableQuote)
			}
			last[q.Pair][q.Exchange] = q

			// Рассмотрим для данной пары все биржи: найдём минимальный ask и максимальный bid
			var (
				bestBuyEx  string
				bestBuy    float64 = math.Inf(1)
				bestSellEx string
				bestSell   float64 = math.Inf(-1)
			)
			for ex, qq := range last[q.Pair] {
				if qq.Ask > 0 && qq.Ask < bestBuy {
					bestBuy, bestBuyEx = qq.Ask, ex
				}
				if qq.Bid > bestSell {
					bestSell, bestSellEx = qq.Bid, ex
				}
			}

			if bestBuyEx == "" || bestSellEx == "" || bestBuyEx == bestSellEx {
				continue
			}

			// Учёт комиссий: используем taker на покупке и продаже.
			buyTaker := 0.0
			sellTaker := 0.0
			if a, ok := u.adapters[bestBuyEx]; ok {
				_, t := a.TradingFee(q.Pair)
				buyTaker = t
			}
			if a, ok := u.adapters[bestSellEx]; ok {
				_, t := a.TradingFee(q.Pair)
				sellTaker = t
			}

			effBuy := bestBuy * (1 + buyTaker)
			effSell := bestSell * (1 - sellTaker)
			gross := bestSell - bestBuy
			net := effSell - effBuy
			if net <= 0 {
				continue
			}

			opp := entity.ArbOpportunity{
				Pair:       q.Pair,
				BuyOn:      bestBuyEx,
				BuyPrice:   bestBuy,
				SellOn:     bestSellEx,
				SellPrice:  bestSell,
				GrossPnl:   gross,
				NetPnl:     net,
				SpreadPct:  (net / bestBuy) * 100,
				DetectedAt: time.Now(),
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			case out <- opp:
			}
		}
	}
}
