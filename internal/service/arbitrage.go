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
	"github.com/dimryb/cross-arb/internal/storage"
)

type Arbitrage struct {
	ctx   context.Context
	app   i.Application
	log   i.Logger
	cfg   *config.CrossArbConfig
	store *storage.TickerStore
}

func NewArbitrageService(
	ctx context.Context,
	app i.Application,
	logger i.Logger,
	cfg *config.CrossArbConfig,
	store *storage.TickerStore,
) *Arbitrage {
	return &Arbitrage{
		ctx:   ctx,
		app:   app,
		log:   logger,
		cfg:   cfg,
		store: store,
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
					m.updateStore("mexc", r)
					printTicker(&r.Data)
				}
			}
		}
	}()

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}

func (m *Arbitrage) updateStore(exchange string, r Result) {
	if r.Error == nil {
		m.store.Set(storage.TickerData{
			Symbol:   r.Data.Symbol,
			Exchange: exchange,
			BidPrice: parseFloat(r.Data.BidPrice),
			BidQty:   parseFloat(r.Data.BidQty),
			AskPrice: parseFloat(r.Data.AskPrice),
			AskQty:   parseFloat(r.Data.AskQty),
		})
	}
}

func printTicker(t *BookTicker) {
	fmt.Printf(
		"  [%s] -> покупка: %s (%s) | продажа: %s (%s)\n",
		t.Symbol,
		t.BidPrice, t.BidQty,
		t.AskPrice, t.AskQty,
	)
}

func getTicker(sc *spotlist.SpotClient, results []Result, index int, symbol string) {
	ticker, err := bookTicker(sc, symbol)
	if err != nil {
		results[index] = Result{
			Symbol: symbol,
			Data:   BookTicker{},
			Error:  err,
		}
	}
	results[index] = Result{
		Symbol: symbol,
		Data:   ticker,
		Error:  nil,
	}
}

func bookTicker(sc *spotlist.SpotClient, symbol string) (BookTicker, error) {
	params := fmt.Sprintf(`{"symbol":"%s"}`, symbol)
	resp, err := sc.BookTicker(params)
	if err != nil {
		return BookTicker{}, fmt.Errorf("BookTicker request failed: %w", err)
	}

	var tickerData BookTicker
	err = json.Unmarshal(resp.Body(), &tickerData)
	if err != nil {
		return BookTicker{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return tickerData, nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
