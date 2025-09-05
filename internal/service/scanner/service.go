package scanner

import (
	"context"
	"errors"
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

	dexUC *uc.DEXPriceUseCase
	cexUC *uc.CEXOrderBookUseCase
	oppUC *uc.ArbOpportunityUseCase
}

// NewService — принимает готовые зависимости (usecase’ы, каналы, адаптеры).
// Требования: interval>0, baseAmount>0, есть пары и хотя бы один адаптер;
// хотя бы один из выходных каналов не nil; для детектора нужен oppCh.
func NewService(
	log i.Logger,
	interval time.Duration,
	baseAmount float64,
	pairs []string,
	adapters []i.EXAdapter,
	pricesCh chan<- entity.ExecutableQuote,
	orderBooksCh chan<- entity.OrderBookResult,
	oppCh chan<- entity.ArbOpportunity,
	dexUC *uc.DEXPriceUseCase,
	cexUC *uc.CEXOrderBookUseCase,
	oppUC *uc.ArbOpportunityUseCase,
) (*Service, error) {
	if interval <= 0 {
		return nil, errors.New("interval must be > 0")
	}
	if baseAmount <= 0 {
		return nil, errors.New("baseAmount must be > 0")
	}
	if len(pairs) == 0 {
		return nil, errors.New("no pairs provided")
	}
	if len(adapters) == 0 {
		return nil, errors.New("no adapters provided")
	}
	if pricesCh == nil && orderBooksCh == nil && oppCh == nil {
		return nil, errors.New("no output channels provided")
	}
	if dexUC == nil {
		dexUC = uc.NewDEXPriceUseCase()
	}
	if cexUC == nil {
		cexUC = uc.NewCEXOrderBookUseCase(5) // дефолтный limit
	}
	if oppUC == nil {
		oppUC = uc.NewOpportunityUseCase(nil)
	}

	return &Service{
		log:          log,
		interval:     interval,
		baseAmount:   baseAmount,
		pairs:        pairs,
		adapters:     adapters,
		pricesCh:     pricesCh,
		orderBooksCh: orderBooksCh,
		oppCh:        oppCh,
		dexUC:        dexUC,
		cexUC:        cexUC,
		oppUC:        oppUC,
	}, nil
}

// Start запускает DEX-цены, CEX-стаканы и детектор возможностей; блокирует до отмены контекста.
func (s *Service) Start(ctx context.Context) error {
	if ctx == nil {
		return errors.New("nil context")
	}

	dexAdapters, cexAdapters := splitAdapters(s.adapters)
	if len(dexAdapters) == 0 && len(cexAdapters) == 0 {
		return errors.New("no supported adapters (DEX/CEX)")
	}

	// Локальный канал DEX-котировок: потом сделаем fan-out наружу и в детектор
	var dexOut chan entity.ExecutableQuote
	if len(dexAdapters) > 0 {
		dexOut = make(chan entity.ExecutableQuote, 128)
		go func() {
			_ = s.dexUC.Stream(ctx, dexAdapters, s.pairs, s.interval, s.baseAmount, dexOut)
		}()
	}

	// Fan-out: наружу (pricesCh) и в opp-детектор
	var oppIn <-chan entity.ExecutableQuote
	if dexOut != nil {
		if s.pricesCh == nil {
			oppIn = dexOut
		} else {
			fan := make(chan entity.ExecutableQuote, 128)
			oppIn = fan
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case q, ok := <-dexOut:
						if !ok {
							return
						}
						select {
						case <-ctx.Done():
							return
						case s.pricesCh <- q:
						}
						select {
						case <-ctx.Done():
							return
						case fan <- q:
						}
					}
				}
			}()
		}
	}

	// CEX-стаканы напрямую наружу
	if s.orderBooksCh != nil && len(cexAdapters) > 0 {
		go func() {
			_ = s.cexUC.Stream(ctx, cexAdapters, s.pairs, s.interval, s.orderBooksCh)
		}()
	}

	// Детектор возможностей (если есть вход и выход)
	if s.oppCh != nil && oppIn != nil {
		go func() {
			if err := s.oppUC.Detect(ctx, oppIn, s.oppCh); err != nil && s.log != nil {
				s.log.Errorf("opportunity detector stopped: %v", err)
			}
		}()
	}

	<-ctx.Done()
	return ctx.Err()
}

// splitAdapters делит общий список на DEX и CEX.
func splitAdapters(adapters []i.EXAdapter) (dex []i.DEXAdapter, cex []i.CEXAdapter) {
	for _, ad := range adapters {
		if ad == nil {
			continue
		}
		if da, ok := any(ad).(i.DEXAdapter); ok {
			dex = append(dex, da)
		}
		if ca, ok := any(ad).(i.CEXAdapter); ok {
			cex = append(cex, ca)
		}
	}
	return
}
