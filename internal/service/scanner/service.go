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

	dexOut := s.startDEX(ctx, dexAdapters)
	oppIn := s.setupDEXFanOut(ctx, dexOut)
	s.startCEX(ctx, cexAdapters)
	s.startOpportunityDetector(ctx, oppIn)

	<-ctx.Done()
	return ctx.Err()
}

// startDEX запускает поток получения DEX-котировок, если есть DEX-адаптеры.
// Возвращает канал с котировками или nil, если DEX не используется.
func (s *Service) startDEX(ctx context.Context, dexAdapters []i.DEXAdapter) <-chan entity.ExecutableQuote {
	if len(dexAdapters) == 0 {
		return nil
	}

	dexOut := make(chan entity.ExecutableQuote, 128)
	go func() {
		err := s.dexUC.Stream(ctx, dexAdapters, s.pairs, s.interval, s.baseAmount, dexOut)
		if err != nil && s.log != nil {
			s.log.Errorf("DEX stream stopped: %v", err)
		}
		close(dexOut)
	}()
	return dexOut
}

// setupDEXFanOut организует fan-out из DEX-канала: в pricesCh (если задан) и в отдельный канал для детектора.
// Возвращает канал для детектора (oppIn) или nil, если входных данных нет.
func (s *Service) setupDEXFanOut(
	ctx context.Context,
	dexOut <-chan entity.ExecutableQuote,
) <-chan entity.ExecutableQuote {
	if dexOut == nil {
		return nil
	}

	if s.pricesCh == nil {
		return dexOut // Нет внешнего канала — детектор получает напрямую
	}

	fan := make(chan entity.ExecutableQuote, 128)
	go func() {
		defer close(fan)
		for {
			select {
			case <-ctx.Done():
				return
			case q, ok := <-dexOut:
				if !ok {
					return
				}
				// Отправка во внешний канал
				select {
				case <-ctx.Done():
					return
				case s.pricesCh <- q:
				}
				// Отправка в fan-out для детектора
				select {
				case <-ctx.Done():
					return
				case fan <- q:
				}
			}
		}
	}()
	return fan
}

// startCEX запускает поток получения CEX-стаканов, если есть CEX-адаптеры и выходной канал.
func (s *Service) startCEX(ctx context.Context, cexAdapters []i.CEXAdapter) {
	if s.orderBooksCh == nil || len(cexAdapters) == 0 {
		return
	}

	go func() {
		err := s.cexUC.Stream(ctx, cexAdapters, s.pairs, s.interval, s.orderBooksCh)
		if err != nil && s.log != nil {
			s.log.Errorf("CEX stream stopped: %v", err)
		}
	}()
}

// startOpportunityDetector запускает детектор арбитражных возможностей, если есть вход и выход.
func (s *Service) startOpportunityDetector(ctx context.Context, oppIn <-chan entity.ExecutableQuote) {
	if s.oppCh == nil || oppIn == nil {
		return
	}

	go func() {
		err := s.oppUC.Detect(ctx, oppIn, s.oppCh)
		if err != nil && s.log != nil {
			s.log.Errorf("opportunity detector stopped: %v", err)
		}
	}()
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
