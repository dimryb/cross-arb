package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	spotlist "github.com/dimryb/cross-arb/internal/api/mexc/spot"
	"github.com/dimryb/cross-arb/internal/api/mexc/utils"
	"github.com/dimryb/cross-arb/internal/config"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/storage"
)

const (
	mexcExchange = "mexc"
	jupExchange  = "jupiter"
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

	mexcCfg, ok := m.cfg.Exchanges[mexcExchange]
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
						getMexcTicker(spot, results, ind, symbol)
					}()
				}
				wgSymbols.Wait()

				m.updateAllStores(mexcExchange, results)
				printTickersReport(results)
			}
		}
	}()

	err := m.runJupiterClient(wg)
	if err != nil {
		return err
	}

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}

func (m *Arbitrage) updateAllStores(exchange string, results []Result) {
	for _, r := range results {
		m.updateStore(exchange, r)
	}
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

func printTickersReport(results []Result) {
	fmt.Printf("=== Обновление цен (%s) ===\n", time.Now().Format("15:04:05.000"))
	for _, r := range results {
		printTicker(r.Data)
	}
}

func printTicker(t BookTicker) {
	fmt.Printf(
		"  [%s] -> покупка: %s (%s) | продажа: %s (%s)\n",
		t.Symbol,
		t.BidPrice, t.BidQty,
		t.AskPrice, t.AskQty,
	)
}

func (m *Arbitrage) runJupiterClient(wg *sync.WaitGroup) error {
	jupiterCfg, ok := m.cfg.Exchanges[jupExchange]
	if !ok || !jupiterCfg.Enabled {
		return fmt.Errorf("jupiter exchange not configured or disabled")
	}
	jupClient, err := jupiter.NewJupiterClient(m.log, jupiterCfg.BaseURL)
	if err != nil {
		return fmt.Errorf("failed to init jupiter client: %w", err)
	}

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
						bookTicker, err := getJupiterTicker(jupClient, symbol)
						processTickerResult(results, ind, symbol, bookTicker, err)
					}()
				}

				wgSymbols.Wait()

				m.updateAllStores(jupExchange, results)
				printTickersReport(results)
			}
		}
	}()
	return nil
}

// getJupiterTicker запрашивает котировку Jupiter и преобразует её в BookTicker.
func getJupiterTicker(jc *jupiter.Client, symbol string) (BookTicker, error) {
	inMint, outMint, err := jupiter.ConvertSpotToMints(symbol)
	if err != nil {
		return BookTicker{}, fmt.Errorf("unsupported symbol format %q", symbol)
	}

	base, quote, err := jupiter.ParseSpotSymbol(symbol)
	if err != nil {
		return BookTicker{}, err
	}

	// Получаем единичные количества для нормализации
	baseUnit, err := jupiter.UnitAmount(base)
	if err != nil {
		return BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", base, err)
	}
	quoteUnit, err := jupiter.UnitAmount(quote)
	if err != nil {
		return BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", quote, err)
	}

	// Запрос 1: base → quote (ASK - цена продажи базового актива)
	askQuote, err := jc.Quote(context.Background(), inMint, outMint, baseUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return BookTicker{}, fmt.Errorf("failed to get ask quote: %w", err)
	}

	// Запрос 2: quote → base (BID - сколько базового актива получим за единицу котировочного)
	bidQuote, err := jc.Quote(context.Background(), outMint, inMint, quoteUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return BookTicker{}, fmt.Errorf("failed to get bid quote: %w", err)
	}

	askPrice := calculatePrice(askQuote.InAmount, askQuote.OutAmount, baseUnit, quoteUnit, false)
	bidPrice := calculatePrice(bidQuote.InAmount, bidQuote.OutAmount, quoteUnit, baseUnit, true)

	return BookTicker{
		Symbol:   symbol,
		BidPrice: fmt.Sprintf("%.6f", bidPrice),
		BidQty:   "0",
		AskPrice: fmt.Sprintf("%.6f", askPrice),
		AskQty:   "0",
	}, nil
}

func calculatePrice(inAmount, outAmount string, inUnit, outUnit int64, invert bool) float64 {
	inAmt := parseFloat(inAmount)
	outAmt := parseFloat(outAmount)

	if inAmt == 0 || outAmt == 0 {
		return 0.0
	}

	inReal := inAmt / float64(inUnit)
	outReal := outAmt / float64(outUnit)

	if invert {
		return inReal / outReal
	}
	return outReal / inReal
}

func getMexcTicker(sc *spotlist.SpotClient, results []Result, index int, symbol string) {
	ticker, err := bookMexcTicker(sc, symbol)
	processTickerResult(results, index, symbol, ticker, err)
}

func processTickerResult(results []Result, index int, symbol string, ticker BookTicker, err error) {
	if err != nil {
		results[index] = Result{
			Symbol: symbol,
			Data:   BookTicker{},
			Error:  err,
		}
	} else {
		results[index] = Result{
			Symbol: symbol,
			Data:   ticker,
			Error:  nil,
		}
	}
}

func bookMexcTicker(sc *spotlist.SpotClient, symbol string) (BookTicker, error) {
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
