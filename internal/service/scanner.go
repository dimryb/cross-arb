package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	i "github.com/dimryb/cross-arb/internal/interface"
)

// Scanner опрашивает набор бирж по списку торговых пар и
// публикует арбитражные возможности подписчикам.
type Scanner struct {
	logger   i.Logger
	interval time.Duration
	pairs    []string
	adapters []ExchangeAdapter

	mu   sync.RWMutex
	subs map[string][]chan Opportunity
}

// Option — функциональная опция инициализации Scanner.
type Option func(*Scanner) error

// WithInterval задаёт период сканирования.
func WithInterval(d time.Duration) Option {
	return func(s *Scanner) error {
		if d <= 0 {
			return fmt.Errorf("интервал должен быть > 0, получено %s", d)
		}
		s.interval = d
		return nil
	}
}

// WithPairs задаёт список торговых пар.
func WithPairs(pairs ...string) Option {
	return func(s *Scanner) error {
		if len(pairs) == 0 {
			return errors.New("нужна хотя бы одна торговая пара")
		}
		s.pairs = append([]string(nil), pairs...) // копируем
		return nil
	}
}

// WithAdapters регистрирует биржевые адаптеры.
func WithAdapters(adapters ...ExchangeAdapter) Option {
	return func(s *Scanner) error {
		if len(adapters) < 2 {
			return errors.New("для арбитража нужны минимум две биржи")
		}
		s.adapters = append([]ExchangeAdapter(nil), adapters...)
		return nil
	}
}

// NewScanner создаёт настроенный сканер.
func NewScanner(l i.Logger, opts ...Option) *Scanner {
	s := &Scanner{
		logger:   l.Named("scanner"),
		interval: 3 * time.Second,
		subs:     make(map[string][]chan Opportunity),
	}
	for _, o := range opts {
		if err := o(s); err != nil {
			l.Fatal("некорректная опция сканера", slog.Any("err", err))
		}
	}
	if len(s.pairs) == 0 {
		l.Fatal("не заданы торговые пары для сканера")
	}
	if len(s.adapters) < 2 {
		l.Fatal("необходимо минимум два адаптера")
	}
	return s
}

// Subscribe возвращает канал Opportunities по паре.
// buf = 0 создаёт небуферизированный канал.
func (s *Scanner) Subscribe(pair string, buf int) (<-chan Opportunity, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if buf < 0 {
		return nil, fmt.Errorf("некорректный размер буфера: %d", buf)
	}
	ch := make(chan Opportunity, buf)
	s.subs[pair] = append(s.subs[pair], ch)
	return ch, nil
}

// Run запускает бесконечный цикл сканирования.
func (s *Scanner) Run(ctx context.Context) error {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case t := <-ticker.C:
			if err := s.scanOnce(ctx, t); err != nil {
				s.logger.Warn("ошибка сканирования", slog.Any("err", err))
			}
		}
	}
}

// scanOnce выполняет один проход по всем биржам для каждой пары.
func (s *Scanner) scanOnce(ctx context.Context, now time.Time) error {
	for _, pair := range s.pairs {
		var wg sync.WaitGroup
		quotesCh := make(chan PricePoint, len(s.adapters))
		errCh := make(chan error, len(s.adapters))

		for _, adp := range s.adapters {
			wg.Add(1)
			go func(a ExchangeAdapter) {
				defer wg.Done()
				bid, ask, err := a.OrderBookTop(ctx, pair)
				if err != nil {
					errCh <- fmt.Errorf("%s: %w", a.Name(), err)
					return
				}
				quotesCh <- PricePoint{
					Exchange:  a.Name(),
					Pair:      pair,
					Bid:       bid,
					Ask:       ask,
					Timestamp: now,
				}
			}(adp)
		}

		wg.Wait()
		close(quotesCh)
		close(errCh)

		// если есть хотя бы одна ошибка — возвращаем первую
		for err := range errCh {
			s.logger.Error("не удалось получить котировки",
				slog.String("pair", pair), slog.Any("err", err))
			return err
		}

		quotes := make([]PricePoint, 0, len(s.adapters))
		for q := range quotesCh {
			quotes = append(quotes, q)
		}

		s.logger.Debug("quotes", slog.String("pair", pair), slog.Any("quotes", quotes))

		if len(quotes) < 2 {
			continue
		}

		bestBidIdx, bestAskIdx := 0, 0
		for i := 1; i < len(quotes); i++ {
			if quotes[i].Bid > quotes[bestBidIdx].Bid {
				bestBidIdx = i
			}
			if quotes[i].Ask < quotes[bestAskIdx].Ask {
				bestAskIdx = i
			}
		}
		if bestBidIdx == bestAskIdx {
			continue
		}

		buy := quotes[bestAskIdx]
		sell := quotes[bestBidIdx]

		_, buyTaker := s.adapters[bestAskIdx].TradingFee(pair)
		_, sellTaker := s.adapters[bestBidIdx].TradingFee(pair)

		gross := sell.Bid - buy.Ask
		net := gross - (buy.Ask*buyTaker + sell.Bid*sellTaker)

		if net <= 0 {
			s.logger.Debug("spread not profitable",
				slog.String("pair", pair),
				slog.Float64("buy_ask", buy.Ask),
				slog.Float64("sell_bid", sell.Bid),
				slog.Float64("net", net),
			)
			continue
		}

		opp := Opportunity{
			Pair:       pair,
			BuyOn:      buy.Exchange,
			BuyPrice:   buy.Ask,
			SellOn:     sell.Exchange,
			SellPrice:  sell.Bid,
			GrossPnl:   gross,
			NetPnl:     net,
			SpreadPct:  net / buy.Ask * 100,
			DetectedAt: now,
		}

		s.logger.Info("обнаружен арбитраж",
			slog.String("pair", opp.Pair),
			slog.String("buy_on", opp.BuyOn),
			slog.Float64("buy_price", opp.BuyPrice),
			slog.String("sell_on", opp.SellOn),
			slog.Float64("sell_price", opp.SellPrice),
			slog.Float64("net", opp.NetPnl),
			slog.Float64("spread_pct", opp.SpreadPct),
		)

		s.publish(opp)
	}

	return nil
}

// publish рассылает Opportunity всем подписчикам пары.
func (s *Scanner) publish(opp Opportunity) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, ch := range s.subs[opp.Pair] {
		select {
		case ch <- opp:
		default:
			s.logger.Warn("канал подписчика переполнен; событие отброшено")
		}
	}
}
