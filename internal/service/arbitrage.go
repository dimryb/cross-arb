package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	spotlist "github.com/dimryb/cross-arb/internal/api/mexc/spot"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
)

type Arbitrage struct {
	ctx context.Context
	app i.Application
	log i.Logger
	cfg *config.CrossArbConfig
}

func NewArbitrageService(
	ctx context.Context,
	app i.Application,
	logger i.Logger,
	cfg *config.CrossArbConfig,
) *Arbitrage {
	return &Arbitrage{
		ctx: ctx,
		app: app,
		log: logger,
		cfg: cfg,
	}
}

type Result struct {
	Symbol string
	Data   interface{}
	Error  error
}

func (m *Arbitrage) Run() error {
	wg := &sync.WaitGroup{}

	mexcCfg, ok := m.cfg.Exchanges["mexc"]
	if !ok || !mexcCfg.Enabled {
		return fmt.Errorf("mexc exchange not configured")
	}
	spot := spotlist.NewSpotClient(m.log, mexcCfg.BaseURL)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				results := make([]Result, len(m.cfg.Symbols))
				wgSymbols := &sync.WaitGroup{}
				for ind, symbol := range m.cfg.Symbols {
					wgSymbols.Add(1)
					go func(index int, sym string) {
						defer wgSymbols.Done()

						params := fmt.Sprintf(`{"symbol":"%s"}`, sym)
						data := spot.BookTicker(params)
						results[index] = Result{
							Symbol: sym,
							Data:   data,
							Error:  nil, // TODO: можно улучшить: если BookTicker возвращает err
						}
					}(ind, symbol)
				}

				wgSymbols.Wait()

				fmt.Printf("=== Обновление цен (%s) ===\n", time.Now().Format("15:04:05.000"))
				for _, r := range results {
					fmt.Printf("  [%s] -> %v\n", r.Symbol, r.Data)
				}
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}
