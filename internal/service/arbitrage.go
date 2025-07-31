package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	spotlist "github.com/dimryb/cross-arb/internal/api/mexc/spot"
	"github.com/dimryb/cross-arb/internal/api/mexc/utils"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/report"
	"github.com/dimryb/cross-arb/internal/types"
)

const (
	exchange = "mexc"
)

type Arbitrage struct {
	ctx   context.Context
	app   i.Application
	log   i.Logger
	cfg   *config.CrossArbConfig
	store i.TickerStore
}

func NewArbitrageService(
	app i.Application,
	cfg *config.CrossArbConfig,
) *Arbitrage {
	return &Arbitrage{
		ctx:   app.Context(),
		app:   app,
		log:   app.Logger(),
		cfg:   cfg,
		store: app.TickerStore(),
	}
}

func (m *Arbitrage) Run() error {
	wg := &sync.WaitGroup{}

	mexcCfg, ok := m.cfg.Exchanges[exchange]
	if !ok || !mexcCfg.Enabled {
		return fmt.Errorf("mexc exchange not configured")
	}
	client := utils.NewClient(mexcCfg.APIKey, mexcCfg.SecretKey, m.log)
	spot := spotlist.NewSpotClient(m.log, mexcCfg.BaseURL, client)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				results := make([]types.Result, len(m.cfg.Symbols))
				wgSymbols := &sync.WaitGroup{}
				for ind, symbol := range m.cfg.Symbols {
					wgSymbols.Add(1)
					go func() {
						defer wgSymbols.Done()
						getTicker(spot, results, ind, symbol)
					}()
				}
				wgSymbols.Wait()

				m.updateAllStores(exchange, results)
				report.PrintTickersReport(results)
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}

func (m *Arbitrage) updateAllStores(exchange string, results []types.Result) {
	for _, r := range results {
		m.updateStore(exchange, r)
	}
}

func (m *Arbitrage) updateStore(exchange string, r types.Result) {
	if r.Error == nil {
		m.store.Set(types.TickerData{
			Symbol:   r.Data.Symbol,
			Exchange: exchange,
			BidPrice: parseFloat(r.Data.BidPrice),
			BidQty:   parseFloat(r.Data.BidQty),
			AskPrice: parseFloat(r.Data.AskPrice),
			AskQty:   parseFloat(r.Data.AskQty),
		})
	}
}

func getTicker(sc *spotlist.SpotClient, results []types.Result, index int, symbol string) {
	ticker, err := bookTicker(sc, symbol)
	if err != nil {
		results[index] = types.Result{
			Symbol: symbol,
			Data:   types.BookTicker{},
			Error:  err,
		}
	}
	results[index] = types.Result{
		Symbol: symbol,
		Data:   ticker,
		Error:  nil,
	}
}

func bookTicker(sc *spotlist.SpotClient, symbol string) (types.BookTicker, error) {
	params := fmt.Sprintf(`{"symbol":"%s"}`, symbol)
	resp, err := sc.BookTicker(params)
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("BookTicker request failed: %w", err)
	}

	var tickerData types.BookTicker
	err = json.Unmarshal(resp.Body(), &tickerData)
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return tickerData, nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
