package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

// MexcAdapter реализует доступ к публичному REST-API биржи MEXC.
// Используется только энд-поинт depth, поэтому ключ и секрет
// не обязательны. Для боевой торговли стоит добавить WebSocket-стримы.
type MexcAdapter struct {
	client  *http.Client
	baseURL string
	logger  *zap.Logger
}

// NewMexcAdapter возвращает готовый к работе адаптер.
// clientTimeout — таймаут HTTP-запросов; при 0 берётся 3 сек.
func NewMexcAdapter(l *zap.Logger, clientTimeout time.Duration) *MexcAdapter {
	if clientTimeout <= 0 {
		clientTimeout = 3 * time.Second
	}
	return &MexcAdapter{
		client:  &http.Client{Timeout: clientTimeout},
		baseURL: "https://api.mexc.com",
		logger:  l.Named("mexc"),
	}
}

// Name удовлетворяет интерфейсу ExchangeAdapter.
func (m *MexcAdapter) Name() string { return "mexc" }

// OrderBookTop запрашивает топ стакана:
//
//	GET /api/v3/depth?symbol=<SYMBOL>&limit=5
//
// SYMBOL формируем из пары, убирая «/» и приводя к верхнему регистру.
func (m *MexcAdapter) OrderBookTop(ctx context.Context, pair string) (bestBid, bestAsk float64, err error) {
	symbol := strings.ReplaceAll(strings.ToUpper(pair), "/", "")
	url := fmt.Sprintf("%s/api/v3/depth?symbol=%s&limit=5", m.baseURL, symbol)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("создание запроса: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return 0, 0, fmt.Errorf("выполнение запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("код ответа %d", resp.StatusCode)
	}

	var raw struct {
		Bids [][]string `json:"bids"`
		Asks [][]string `json:"asks"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return 0, 0, fmt.Errorf("декодирование JSON: %w", err)
	}
	if len(raw.Bids) == 0 || len(raw.Asks) == 0 {
		return 0, 0, fmt.Errorf("пустой стакан для %s", pair)
	}

	if _, err := fmt.Sscan(raw.Bids[0][0], &bestBid); err != nil {
		return 0, 0, fmt.Errorf("парсинг bid: %w", err)
	}
	if _, err := fmt.Sscan(raw.Asks[0][0], &bestAsk); err != nil {
		return 0, 0, fmt.Errorf("парсинг ask: %w", err)
	}
	return bestBid, bestAsk, nil
}

// TradingFee возвращает фиксированную комиссию MEXC для спота: 0.1 %.
func (m *MexcAdapter) TradingFee(string) (maker, taker float64) { return 0.001, 0.001 }

// Close удовлетворяет интерфейсу, доп. ресурсы не удерживаются.
func (m *MexcAdapter) Close() error { return nil }
