package service

import (
	"context"
	"encoding/json"
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
	Data   BookTicker
	Error  error
}

type BookTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
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
					go func() {
						defer wgSymbols.Done()
						getTicker(spot, results, ind, symbol)
					}()
				}
				wgSymbols.Wait()

				fmt.Printf("=== Обновление цен (%s) ===\n", time.Now().Format("15:04:05.000"))
				for _, r := range results {
					fmt.Printf(
						"  [%s] -> покупка: %s (%s) продажа: %s (%s)\n",
						r.Data.Symbol,
						r.Data.BidPrice, r.Data.BidQty,
						r.Data.AskPrice, r.Data.AskQty,
					)
				}
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}

func getTicker(sc *spotlist.SpotClient, results []Result, index int, sym string) {
	params := fmt.Sprintf(`{"symbol":"%s"}`, sym)
	resp := sc.BookTicker(params)

	var tickerData BookTicker
	err := json.Unmarshal(resp.Body(), &tickerData)
	if err != nil {
		results[index] = Result{
			Symbol: sym,
			Data:   BookTicker{},
			Error:  fmt.Errorf("failed to parse JSON: %w", err),
		}
		return
	}

	results[index] = Result{
		Symbol: sym,
		Data:   tickerData,
		Error:  nil, // TODO: можно улучшить: если BookTicker возвращает err
	}
}
