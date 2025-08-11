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
	"github.com/dimryb/cross-arb/internal/report"
	"github.com/dimryb/cross-arb/internal/types"
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
				results := make([]types.Result, len(m.cfg.Symbols))
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
			}
		}
	}()

	m.runMexcOrderBook(wg)

	err := m.runJupiterClient(wg)
	if err != nil {
		return err
	}

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
				results := make([]types.Result, len(m.cfg.Symbols))
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
			}
		}
	}()
	return nil
}

// getJupiterTicker запрашивает котировку Jupiter и преобразует её в BookTicker.
func getJupiterTicker(jc *jupiter.Client, symbol string) (types.BookTicker, error) {
	inMint, outMint, err := jupiter.ConvertSpotToMints(symbol)
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("unsupported symbol format %q", symbol)
	}

	base, quote, err := jupiter.ParseSpotSymbol(symbol)
	if err != nil {
		return types.BookTicker{}, err
	}

	// Получаем единичные количества для нормализации
	baseUnit, err := jupiter.UnitAmount(base)
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", base, err)
	}
	quoteUnit, err := jupiter.UnitAmount(quote)
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", quote, err)
	}

	// Запрос 1: base → quote (ASK - цена продажи базового актива)
	askQuote, err := jc.Quote(context.Background(), inMint, outMint, baseUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("failed to get ask quote: %w", err)
	}

	// Запрос 2: quote → base (BID - сколько базового актива получим за единицу котировочного)
	bidQuote, err := jc.Quote(context.Background(), outMint, inMint, quoteUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return types.BookTicker{}, fmt.Errorf("failed to get bid quote: %w", err)
	}

	askPrice := calculatePrice(askQuote.InAmount, askQuote.OutAmount, baseUnit, quoteUnit, false)
	bidPrice := calculatePrice(bidQuote.InAmount, bidQuote.OutAmount, quoteUnit, baseUnit, true)

	return types.BookTicker{
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

func getMexcTicker(sc *spotlist.SpotClient, results []types.Result, index int, symbol string) {
	ticker, err := bookMexcTicker(sc, symbol)
	processTickerResult(results, index, symbol, ticker, err)
}

func processTickerResult(results []types.Result, index int, symbol string, ticker types.BookTicker, err error) {
	if err != nil {
		results[index] = types.Result{
			Symbol: symbol,
			Data:   types.BookTicker{},
			Error:  err,
		}
	} else {
		results[index] = types.Result{
			Symbol: symbol,
			Data:   ticker,
			Error:  nil,
		}
	}
}

func bookMexcTicker(sc *spotlist.SpotClient, symbol string) (types.BookTicker, error) {
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

func (m *Arbitrage) runMexcOrderBook(wg *sync.WaitGroup) {
	mexcCfg, ok := m.cfg.Exchanges[mexcExchange]
	if !ok || !mexcCfg.Enabled {
		m.log.Warnf("MEXC exchange not enabled for order book")
		return
	}

	client := utils.NewClient(mexcCfg.APIKey, mexcCfg.SecretKey, m.log)
	spot := spotlist.NewSpotClient(m.log, mexcCfg.BaseURL, client)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				results := make([]types.OrderBookResult, len(m.cfg.Symbols))
				wgSymbols := &sync.WaitGroup{}
				for ind, symbol := range m.cfg.Symbols {
					wgSymbols.Add(1)
					go func(index int, sym string) {
						defer wgSymbols.Done()
						getMexcOrder(spot, results, index, sym, mexcCfg.OrderLimit)
					}(ind, symbol)
				}
				wgSymbols.Wait()

				report.PrintOrderBookReport(results)
				m.findBestOrder(results, mexcExchange, mexcCfg.MaxPriceDiff, mexcCfg.MinQtyImprovement)
			}
		}
	}()
}

func (m *Arbitrage) findBestOrder(
	results []types.OrderBookResult,
	exchange string,
	maxPriceDiff float64,
	minQtyImprovement float64,
) {
	for _, r := range results {
		if r.Error != nil {
			fmt.Printf("  [%s] Error: %v\n", r.Symbol, r.Error)
			continue
		}

		var bestBidPrice, bestBidQty float64
		topBidPrice := r.Data.Bids[0].Price
		topBidQty := r.Data.Bids[0].Quantity

		bestBidPrice = topBidPrice
		bestBidQty = topBidQty

		// Различия в цене для продажи так же учитывается
		for _, bid := range r.Data.Bids {
			if bid.Price < topBidPrice-maxPriceDiff-1e-8 {
				break
			}
			if bid.Quantity >= topBidQty+minQtyImprovement {
				bestBidPrice = bid.Price
				bestBidQty = bid.Quantity
			}
		}

		var bestAskPrice, bestAskQty float64
		topAskPrice := r.Data.Asks[0].Price
		topAskQty := r.Data.Asks[0].Quantity

		bestAskPrice = topAskPrice
		bestAskQty = topAskQty

		for _, ask := range r.Data.Asks {
			if ask.Price > topAskPrice+maxPriceDiff+1e-8 {
				break
			}
			if ask.Quantity >= topAskQty+minQtyImprovement {
				bestAskPrice = ask.Price
				bestAskQty = ask.Quantity
			}
		}

		fmt.Printf("---\n[%s] ЛУЧШИЙ ВАРИАНТ ИЗ СТАКАНА:\n", r.Symbol)
		fmt.Printf("Покупка: %.2f USDT | Количество: %.3f\n", bestAskPrice, bestAskQty)
		fmt.Printf("Продать: %.2f USDT | Количество: %.3f\n", bestBidPrice, bestBidQty)

		m.store.Set(types.TickerData{
			Symbol:   r.Symbol,
			Exchange: exchange,
			BidPrice: bestBidPrice,
			BidQty:   bestBidQty,
			AskPrice: bestAskPrice,
			AskQty:   bestAskQty,
		})
	}
}

func getMexcOrder(sc *spotlist.SpotClient, results []types.OrderBookResult, index int, symbol string, limit int) {
	book, err := bookMexcOrder(sc, symbol, limit)
	processOrderResult(results, index, symbol, book, err)
}

func processOrderResult(results []types.OrderBookResult, index int, symbol string, book types.OrderBook, err error) {
	if err != nil {
		results[index] = types.OrderBookResult{
			Symbol: symbol,
			Data:   types.OrderBook{},
			Error:  err,
		}
	} else {
		results[index] = types.OrderBookResult{
			Symbol: symbol,
			Data:   book,
			Error:  nil,
		}
	}
}

func bookMexcOrder(sc *spotlist.SpotClient, symbol string, limit int) (types.OrderBook, error) {
	params := fmt.Sprintf(`{"symbol":"%s", "limit":"%d"}`, symbol, limit)
	resp, err := sc.Depth(params)
	if err != nil {
		return types.OrderBook{}, fmt.Errorf("MEXC Depth request failed: %w", err)
	}

	var raw struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}
	err = json.Unmarshal(resp.Body(), &raw)
	if err != nil {
		return types.OrderBook{}, fmt.Errorf("failed to parse MEXC Depth JSON: %w", err)
	}

	var bids, asks []types.Order
	for _, item := range raw.Bids {
		if len(item) != 2 {
			continue
		}
		price := parseFloat(item[0])
		qty := parseFloat(item[1])
		if price > 0 && qty > 0 {
			bids = append(bids, types.Order{Price: price, Quantity: qty})
		}
	}
	for _, item := range raw.Asks {
		if len(item) != 2 {
			continue
		}
		price := parseFloat(item[0])
		qty := parseFloat(item[1])
		if price > 0 && qty > 0 {
			asks = append(asks, types.Order{Price: price, Quantity: qty})
		}
	}

	return types.OrderBook{Bids: bids, Asks: asks}, nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
