package mexc

import (
	"net/http"
	"time"

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
		baseURL: "https://api.mexc.com",
		logger:  l.Named("mexc"),
	}
}

// Name удовлетворяет интерфейсу EXAdapter.
func (m *Adapter) Name() string { return "mexc" }

// TradingFee возвращает фиксированную комиссию MEXC для спота: 0.1 %.
func (m *Adapter) TradingFee(string) (maker, taker float64) { return 0.001, 0.001 }

// Close удовлетворяет интерфейсу, доп. ресурсы не удерживаются.
func (m *Adapter) Close() error { return nil }
