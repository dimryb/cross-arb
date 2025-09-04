package scanner

import (
	"errors"
	"fmt"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
	uc "github.com/dimryb/cross-arb/internal/usecase/scan"
)

// Service — оркестратор юзкейсов и публикации событий через каналы.
// Внешние каналы передаются в NewService и НЕ закрываются сервисом.
type Service struct {
	log        i.Logger
	interval   time.Duration
	baseAmount float64
	pairs      []string
	adapters   []i.EXAdapter

	pricesCh     chan<- entity.ExecutableQuote
	orderBooksCh chan<- entity.OrderBookResult
	oppCh        chan<- entity.ArbOpportunity

	dexUC uc.DEXPriceUseCase
	cexUC uc.CEXOrderBookUseCase
	oppUC uc.ArbOpportunityUseCase
}

// NewService создает сканер (оркестратор) для набора DEX/CEX адаптеров.
// Не запускает горутины и не закрывает внешние каналы.
//
// Коротко — что нужно:
//   - interval > 0
//   - len(pairs) > 0
//   - len(adapters) >= 2
//   - oppCh != nil
//
// Режимы:
//   - Есть DEX: pricesCh != nil и baseAmount > 0
//   - Только CEX: orderBooksCh != nil
//
// Nil-юзкейсы автоматически заменяются на Noop-реализации.
func NewService(
	log i.Logger,
	interval time.Duration,
	baseAmount float64,
	pairs []string,
	adapters []i.EXAdapter,
	pricesCh chan<- entity.ExecutableQuote,
	orderBooksCh chan<- entity.OrderBookResult,
	oppCh chan<- entity.ArbOpportunity,
	dexUC uc.DEXPriceUseCase,
	cexUC uc.CEXOrderBookUseCase,
	oppUC uc.ArbOpportunityUseCase,
) (*Service, error) {
	// Быстрая валидация входных данных
	if len(pairs) == 0 {
		return nil, errors.New("no pairs provided")
	}
	if len(adapters) < 2 {
		return nil, errors.New("must be at least 2 adapters")
	}
	if interval <= 0 {
		return nil, fmt.Errorf("interval must be positive, got %s", interval)
	}
	if oppCh == nil {
		return nil, errors.New("oppCh must not be nil")
	}

	hasDEX, hasCEX := detectAdapterKinds(adapters)
	if !hasDEX && !hasCEX {
		return nil, errors.New("no supported adapters")
	}

	if hasDEX {
		if pricesCh == nil {
			return nil, errors.New("pricesCh must not be nil when DEX adapters are used")
		}
		if baseAmount <= 0 {
			return nil, errors.New("baseAmount must be positive when DEX adapters are used")
		}
	}
	// Если только CEX — нужен orderBooksCh
	if !hasDEX && hasCEX && orderBooksCh == nil {
		return nil, errors.New("orderBooksCh must not be nil when only CEX adapters are used")
	}

	// Значения по умолчанию
	if dexUC == nil {
		dexUC = uc.NewNoopDEXPriceUseCase()
	}
	if cexUC == nil {
		cexUC = uc.NewNoopCEXOrderBookUseCase()
	}
	if oppUC == nil {
		oppUC = uc.NewOpportunityUseCase()
	}

	return &Service{
		log:          log,
		interval:     interval,
		baseAmount:   baseAmount,
		pairs:        append([]string(nil), pairs...),
		adapters:     append([]i.EXAdapter(nil), adapters...),
		pricesCh:     pricesCh,
		orderBooksCh: orderBooksCh,
		oppCh:        oppCh,
		dexUC:        dexUC,
		cexUC:        cexUC,
		oppUC:        oppUC,
	}, nil
}

// detectAdapterKinds reports whether the provided set contains any DEX and/or CEX adapters.
func detectAdapterKinds(adapters []i.EXAdapter) (hasDEX, hasCEX bool) {
	for _, ad := range adapters {
		if !hasDEX {
			if _, ok := any(ad).(i.DEXAdapter); ok {
				hasDEX = true
			}
		}
		if !hasCEX {
			if _, ok := any(ad).(i.CEXAdapter); ok {
				hasCEX = true
			}
		}
		if hasDEX && hasCEX {
			break
		}
	}
	return
}
