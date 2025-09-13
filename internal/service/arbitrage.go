package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	"github.com/dimryb/cross-arb/internal/api/mexc"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/dimryb/cross-arb/internal/report"
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
	ctx context.Context,
	app i.Application,
	logger i.Logger,
	cfg *config.CrossArbConfig,
	store i.TickerStore,
) *Arbitrage {
	return &Arbitrage{
		ctx:   ctx,
		app:   app,
		log:   logger,
		cfg:   cfg,
		store: store,
	}
}

func (m *Arbitrage) Run() error {
	wg := &sync.WaitGroup{}

	mexcCfg, ok := m.cfg.Exchanges[mexcExchange]
	if !ok || !mexcCfg.Enabled {
		return fmt.Errorf("mexc exchange not configured")
	}
	client, err := mexc.NewClient(mexcCfg.APIKey, mexcCfg.SecretKey, mexcCfg.BaseURL, m.log)
	if err != nil {
		return err
	}
	spot := mexc.NewSpotList(m.log, client)

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(time.Second)

		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				results := make([]entity.Result, len(m.cfg.Symbols))
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

	err = m.runJupiterClient(wg)
	if err != nil {
		return err
	}

	m.log.Infof("Arbitrage service is running...")

	wg.Wait()

	return nil
}

func (m *Arbitrage) updateAllStores(exchange string, results []entity.Result) {
	for _, r := range results {
		m.updateStore(exchange, r)
	}
}

func (m *Arbitrage) updateStore(exchange string, r entity.Result) {
	if r.Error == nil {
		m.store.Set(entity.TickerData{
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
				results := make([]entity.Result, len(m.cfg.Symbols))
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
func getJupiterTicker(jc *jupiter.Client, symbol string) (entity.BookTicker, error) {
	inMint, outMint, err := jupiter.ConvertSpotToMints(symbol)
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("unsupported symbol format %q", symbol)
	}

	base, quote, err := jupiter.ParseSpotSymbol(symbol)
	if err != nil {
		return entity.BookTicker{}, err
	}

	// Получаем единичные количества для нормализации
	baseUnit, err := jupiter.UnitAmount(base)
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", base, err)
	}
	quoteUnit, err := jupiter.UnitAmount(quote)
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("failed to get unit amount for %s: %w", quote, err)
	}

	// Запрос 1: base → quote (ASK - цена продажи базового актива)
	askQuote, err := jc.Quote(context.Background(), inMint, outMint, baseUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("failed to get ask quote: %w", err)
	}

	// Запрос 2: quote → base (BID - сколько базового актива получим за единицу котировочного)
	bidQuote, err := jc.Quote(context.Background(), outMint, inMint, quoteUnit, jupiter.DefaultQuoteOptions())
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("failed to get bid quote: %w", err)
	}

	askPrice := calculatePrice(askQuote.InAmount, askQuote.OutAmount, baseUnit, quoteUnit, false)
	bidPrice := calculatePrice(bidQuote.InAmount, bidQuote.OutAmount, quoteUnit, baseUnit, true)

	return entity.BookTicker{
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

func getMexcTicker(sc *mexc.SpotList, results []entity.Result, index int, symbol string) {
	ticker, err := bookMexcTicker(sc, symbol)
	processTickerResult(results, index, symbol, ticker, err)
}

func processTickerResult(results []entity.Result, index int, symbol string, ticker entity.BookTicker, err error) {
	if err != nil {
		results[index] = entity.Result{
			Symbol: symbol,
			Data:   entity.BookTicker{},
			Error:  err,
		}
	} else {
		results[index] = entity.Result{
			Symbol: symbol,
			Data:   ticker,
			Error:  nil,
		}
	}
}

func bookMexcTicker(sc *mexc.SpotList, symbol string) (entity.BookTicker, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := map[string]string{"symbol": symbol}
	resp, err := sc.BookTicker(ctx, params)
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("BookTicker request failed: %w", err)
	}

	var tickerData entity.BookTicker
	err = json.Unmarshal(resp.Body(), &tickerData)
	if err != nil {
		return entity.BookTicker{}, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return tickerData, nil
}

func (m *Arbitrage) runMexcOrderBook(wg *sync.WaitGroup) {
	mexcCfg, ok := m.cfg.Exchanges[mexcExchange]
	if !ok || !mexcCfg.Enabled {
		m.log.Warnf("MEXC exchange not enabled for order book")
		return
	}

	client, err := mexc.NewClient(mexcCfg.APIKey, mexcCfg.SecretKey, mexcCfg.BaseURL, m.log)
	if err != nil {
		m.log.Warnf("MEXC client creation failed: %v", err)
	}
	spot := mexc.NewSpotList(m.log, client)

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
				results := make([]entity.OrderBookResult, len(m.cfg.Symbols))
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
	results []entity.OrderBookResult,
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

		// fmt.Printf("---\n[%s] ЛУЧШИЙ ВАРИАНТ ИЗ СТАКАНА:\n", r.Symbol)
		// fmt.Printf("Покупка: %.2f USDT | Количество: %.3f\n", bestAskPrice, bestAskQty)
		// fmt.Printf("Продать: %.2f USDT | Количество: %.3f\n", bestBidPrice, bestBidQty)

		m.store.Set(entity.TickerData{
			Symbol:   r.Symbol,
			Exchange: exchange,
			BidPrice: bestBidPrice,
			BidQty:   bestBidQty,
			AskPrice: bestAskPrice,
			AskQty:   bestAskQty,
		})
	}
}

func getMexcOrder(sc *mexc.SpotList, results []entity.OrderBookResult, index int, symbol string, limit int) {
	book, err := bookMexcOrder(sc, symbol, limit)
	processOrderResult(results, index, symbol, book, err)
}

func processOrderResult(results []entity.OrderBookResult, index int, symbol string, book entity.OrderBook, err error) {
	if err != nil {
		results[index] = entity.OrderBookResult{
			Symbol: symbol,
			Data:   entity.OrderBook{},
			Error:  err,
		}
	} else {
		results[index] = entity.OrderBookResult{
			Symbol: symbol,
			Data:   book,
			Error:  nil,
		}
	}
}

func bookMexcOrder(sc *mexc.SpotList, symbol string, limit int) (entity.OrderBook, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	params := map[string]string{"symbol": symbol, "limit": fmt.Sprintf("%d", limit)}
	resp, err := sc.Depth(ctx, params)
	if err != nil {
		return entity.OrderBook{}, fmt.Errorf("MEXC Depth request failed: %w", err)
	}

	var raw struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}
	err = json.Unmarshal(resp.Body(), &raw)
	if err != nil {
		return entity.OrderBook{}, fmt.Errorf("failed to parse MEXC Depth JSON: %w", err)
	}

	var bids, asks []entity.Order
	for _, item := range raw.Bids {
		if len(item) != 2 {
			continue
		}
		price := parseFloat(item[0])
		qty := parseFloat(item[1])
		if price > 0 && qty > 0 {
			bids = append(bids, entity.Order{Price: price, Quantity: qty})
		}
	}
	for _, item := range raw.Asks {
		if len(item) != 2 {
			continue
		}
		price := parseFloat(item[0])
		qty := parseFloat(item[1])
		if price > 0 && qty > 0 {
			asks = append(asks, entity.Order{Price: price, Quantity: qty})
		}
	}

	return entity.OrderBook{Bids: bids, Asks: asks}, nil
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}
