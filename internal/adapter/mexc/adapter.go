package mexc

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/cross-arb/internal/api/mexc"
	"github.com/dimryb/cross-arb/internal/config"
	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// Константа для проверки количества полей в одной записи стакана [Price, Quantity].
const OrderBookEntryFields = 2

type AdapterConfig struct {
	APIKey    string
	SecretKey string
	BaseURL   string
	Enabled   bool
	Timeout   time.Duration
	Pairs     map[string]entity.MintPair // symbol → [base_mint, quote_mint]
}

// Adapter реализует доступ к публичному REST-API биржи MEXC.
// Используется только энд-поинт depth, поэтому ключ и секрет
// не обязательны. Для боевой торговли стоит добавить WebSocket-стримы.
type Adapter struct {
	client     mexc.Client
	baseURL    string
	logger     i.Logger
	pairConfig map[string]entity.MintPair // "SOL/USDT" → {baseMint, quoteMint}
}

// NewAdapter возвращает готовый к работе адаптер.
// clientTimeout — таймаут HTTP-запросов; при 0 берётся 3 сек.
func NewAdapter(l i.Logger, cfg *AdapterConfig) *Adapter {
	client, err := mexc.NewMexcClient(cfg.APIKey, cfg.SecretKey, cfg.BaseURL, l)
	if err != nil {
		l.Fatalf("failed to create Mexc client: %v", err)
	}

	return &Adapter{
		client:     *client,
		logger:     l.Named(config.MexcExchange),
		baseURL:    cfg.BaseURL,
		pairConfig: cfg.Pairs,
	}
}

// Name удовлетворяет интерфейсу EXAdapter.
func (m *Adapter) Name() string { return "mexc" }

// TradingFee возвращает фиксированную комиссию MEXC для спота: 0.1 %.
func (m *Adapter) TradingFee(string) (maker, taker float64) { return 0.001, 0.001 }

// Close удовлетворяет интерфейсу, доп. ресурсы не удерживаются.
func (m *Adapter) Close() error { return nil }

// OrderBookDepth возращает стакан.
func (m *Adapter) OrderBookDepth(ctx context.Context, pair string, limit int) (entity.OrderBook, error) {
	// Преобразуем "SOL/USDT" → "SOLUSDT" здесь или в другом месте?
	symbol := strings.ReplaceAll(pair, "/", "")
	params := map[string]string{
		"symbol": symbol,
		"limit":  fmt.Sprintf("%d", limit),
	}

	spot := mexc.NewSpotAPI(m.logger, &m.client)

	depthResp, err := spot.Market.Depth(ctx, params)
	if err != nil {
		return entity.OrderBook{}, fmt.Errorf("ошибка получения ответа от Spot API: %w", err)
	}

	var bids, asks []entity.Order
	for _, item := range depthResp.Bids {
		if len(item) != OrderBookEntryFields {
			continue
		}
		price, err := strconv.ParseFloat(item[0], 64)
		if err != nil {
			return entity.OrderBook{}, fmt.Errorf("ошибка парсинга bids price: %w", err)
		}

		quantity, err := strconv.ParseFloat(item[1], 64)
		if err != nil {
			return entity.OrderBook{}, fmt.Errorf("ошибка парсинга bids количества: %w", err)
		}
		if price > 0 && quantity > 0 {
			bids = append(bids, entity.Order{Price: price, Quantity: quantity})
		}
	}

	for _, item := range depthResp.Asks {
		if len(item) != OrderBookEntryFields {
			continue
		}
		price, err := strconv.ParseFloat(item[0], 64)
		if err != nil {
			return entity.OrderBook{}, fmt.Errorf("ошибка парсинга ask price: %w", err)
		}
		quantity, err := strconv.ParseFloat(item[1], 64)
		if err != nil {
			return entity.OrderBook{}, fmt.Errorf("ошибка парсинга ask quantity: %w", err)
		}
		if price > 0 && quantity > 0 {
			asks = append(asks, entity.Order{Price: price, Quantity: quantity})
		}
	}

	return entity.OrderBook{Bids: bids, Asks: asks}, nil
}
