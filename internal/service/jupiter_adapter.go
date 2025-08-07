package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	i "github.com/dimryb/cross-arb/internal/interface"
)

// JupiterAdapterConfig конфигурация адаптера.
type JupiterAdapterConfig struct {
	BaseURL string
	Enabled bool
	Timeout time.Duration
	Pairs   map[string]MintPair // symbol → [base_mint, quote_mint]
}

// JupiterAdapter использует публичный Quote-API агрегатора Jupiter (Solana).
// Он запрашивает цену обмена base → quote (ask) и quote → base (bid).
type JupiterAdapter struct {
	client     *http.Client
	logger     i.Logger
	baseURL    string
	pairConfig map[string]MintPair // "SOL/USDT" → {baseMint, quoteMint}
}

// NewJupiterAdapter создаёт адаптер.
// pairMap: "SOL/USDT": {baseMint, quoteMint}.
func NewJupiterAdapter(l i.Logger, cfg *JupiterAdapterConfig) *JupiterAdapter {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &JupiterAdapter{
		client:     &http.Client{Timeout: timeout},
		logger:     l.Named("jupiter"),
		baseURL:    cfg.BaseURL,
		pairConfig: cfg.Pairs,
	}
}

// Name удовлетворяет интерфейсу ExchangeAdapter.
func (j *JupiterAdapter) Name() string { return "jupiter" }

// OrderBookTop для Jupiter запрашивает quote двумя направлениями.
func (j *JupiterAdapter) OrderBookTop(ctx context.Context, pair string) (bestBid, bestAsk float64, err error) {
	mints, ok := j.pairConfig[pair]
	if !ok {
		return 0, 0, fmt.Errorf("неизвестная пара %s", pair)
	}

	ask, err := j.quote(ctx, mints.BaseMint, mints.QuoteMint)
	if err != nil {
		return 0, 0, fmt.Errorf("ask: %w", err)
	}
	bid, err := j.quote(ctx, mints.QuoteMint, mints.BaseMint)
	if err != nil {
		return 0, 0, fmt.Errorf("bid: %w", err)
	}
	return bid, ask, nil
}

// quote вызывает /v6/quote для объёма 0.01 токена (1e6 минимальных единиц).
func (j *JupiterAdapter) quote(ctx context.Context, inMint, outMint string) (float64, error) {
	url := fmt.Sprintf(
		"%s?inputMint=%s&outputMint=%s&amount=1000000&slippageBps=1",
		j.baseURL, inMint, outMint,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}
	resp, err := j.client.Do(req)
	if err != nil {
		return 0, err
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("код %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var raw struct {
		Data []struct {
			OutAmount float64 `json:"outAmount,string"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return 0, err
	}
	if len(raw.Data) == 0 {
		return 0, fmt.Errorf("пустой ответ")
	}

	// Цена за 1 единицу base-токена.
	return raw.Data[0].OutAmount / 1e6, nil
}

// TradingFee: Jupiter комиссия 0 (только сеть).
func (j *JupiterAdapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close: дополнительных ресурсов нет.
func (j *JupiterAdapter) Close() error { return nil }
