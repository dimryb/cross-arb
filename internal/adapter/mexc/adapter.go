package mexc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/cross-arb/internal/entity"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// Adapter реализует доступ к публичному REST-API биржи MEXC.
// Используется только энд-поинт depth, поэтому ключ и секрет
// не обязательны. Для боевой торговли стоит добавить WebSocket-стримы.
type Adapter struct {
	client  *http.Client
	baseURL string
	logger  i.Logger
}

// NewAdapter возвращает готовый к работе адаптер.
// clientTimeout — таймаут HTTP-запросов; при 0 берётся 3 сек.
func NewAdapter(l i.Logger, clientTimeout time.Duration) *Adapter {
	if clientTimeout <= 0 {
		clientTimeout = 3 * time.Second
	}
	return &Adapter{
		client:  &http.Client{Timeout: clientTimeout},
		baseURL: "https://api.mexc.com", // базовый URL; конечная точка: /api/v3/depth
		logger:  l.Named("mexc"),
	}
}

// Name удовлетворяет интерфейсу EXAdapter.
func (m *Adapter) Name() string { return "mexc" }

// TradingFee возвращает фиксированную комиссию MEXC для спота: 0.1 %.
func (m *Adapter) TradingFee(string) (maker, taker float64) { return 0.001, 0.001 }

// Close удовлетворяет интерфейсу, доп. ресурсы не удерживаются.
func (m *Adapter) Close() error { return nil }

// OrderBookDepth реализует метод CEXAdapter: возвращает стакан по REST /api/v3/depth.
func (m *Adapter) OrderBookDepth(ctx context.Context, pair string, limit int) (entity.OrderBook, error) {
	var empty entity.OrderBook

	if ctx == nil {
		return empty, fmt.Errorf("nil context")
	}

	symbol := normalizeSymbol(pair) // "SOL/USDT" -> "SOLUSDT"

	if limit <= 0 {
		limit = 5
	}
	// Подстрахуем лимит популярными значениями MEXC (5, 10, 20, 50, 100, 200, 500)
	if limit > 500 {
		limit = 500
	}

	endpoint := fmt.Sprintf("%s/api/v3/depth", strings.TrimRight(m.baseURL, "/"))
	u, err := url.Parse(endpoint)
	if err != nil {
		return empty, fmt.Errorf("parse base url: %w", err)
	}
	q := u.Query()
	q.Set("symbol", symbol)
	q.Set("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return empty, fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return empty, fmt.Errorf("http do: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		var bodyPreview [1024]byte
		n, _ := resp.Body.Read(bodyPreview[:])
		return empty, fmt.Errorf("mexc depth status=%d body=%s", resp.StatusCode, string(bodyPreview[:n]))
	}

	var dr struct {
		LastUpdateID int64      `json:"lastUpdateId"`
		Bids         [][]string `json:"bids"`
		Asks         [][]string `json:"asks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		return empty, fmt.Errorf("decode depth response: %w", err)
	}

	book := entity.OrderBook{
		Bids: make([]entity.Order, 0, len(dr.Bids)),
		Asks: make([]entity.Order, 0, len(dr.Asks)),
	}

	parse := func(pq []string) (entity.Order, error) {
		if len(pq) < 2 {
			return entity.Order{}, fmt.Errorf("unexpected level format: %v", pq)
		}
		price, err := strconv.ParseFloat(pq[0], 64)
		if err != nil {
			return entity.Order{}, fmt.Errorf("parse price: %w", err)
		}
		qty, err := strconv.ParseFloat(pq[1], 64)
		if err != nil {
			return entity.Order{}, fmt.Errorf("parse qty: %w", err)
		}
		return entity.Order{Price: price, Quantity: qty}, nil
	}

	for _, lvl := range dr.Bids {
		if ord, err := parse(lvl); err == nil {
			book.Bids = append(book.Bids, ord)
		} else {
			m.logger.Warn("skip bad bid level", "err", err)
		}
	}
	for _, lvl := range dr.Asks {
		if ord, err := parse(lvl); err == nil {
			book.Asks = append(book.Asks, ord)
		} else {
			m.logger.Warn("skip bad ask level", "err", err)
		}
	}

	return book, nil
}

// normalizeSymbol переводит "BASE/QUOTE" → "BASEQUOTE" и удаляет пробелы.
func normalizeSymbol(pair string) string {
	s := strings.ToUpper(strings.TrimSpace(pair))
	s = strings.ReplaceAll(s, "/", "")
	return s
}
