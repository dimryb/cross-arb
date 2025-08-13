package adapter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
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
	client     *jupiter.Client
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
	//TODO: err do
	client, _ := jupiter.NewJupiterClient(l, cfg.BaseURL)
	return &JupiterAdapter{
		client:     client,
		logger:     l.Named("jupiter"),
		baseURL:    cfg.BaseURL,
		pairConfig: cfg.Pairs,
	}
}

// Name удовлетворяет интерфейсу ExchangeAdapter.
func (j *JupiterAdapter) Name() string { return "jupiter" }

// OrderBookTop для Jupiter: цены в USDT за 1 SOL для пары SOL/USDT.
func (j *JupiterAdapter) OrderBookTop(ctx context.Context, pair string) (bestBid, bestAsk float64, err error) {
	mints, ok := j.pairConfig[pair]
	if !ok {
		return 0, 0, fmt.Errorf("неизвестная пара %s", pair)
	}

	// ask: сколько QUOTE за 1 BASE
	ask, err := j.quote(ctx, mints.BaseMint, mints.QuoteMint, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("ask: %w", err)
	}

	// rawBid: сколько BASE за 1 QUOTE
	rawBid, err := j.quote(ctx, mints.QuoteMint, mints.BaseMint, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("bid: %w", err)
	}
	if rawBid == 0 {
		return 0, 0, fmt.Errorf("zero raw bid")
	}

	// bid в тех же единицах, что и ask: QUOTE за 1 BASE
	bid := 1 / rawBid
	return bid, ask, nil
}

// quote возвращает: "сколько OUT токенов за 1 IN токен".
func (j *JupiterAdapter) quote(ctx context.Context, inMint, outMint string, opts *jupiter.QuoteOptions) (float64, error) {
	inUnit, err := jupiter.UnitAmountByMint(inMint) // 10^decimals(IN)
	if err != nil {
		return 0, err
	}
	resp, err := j.client.Quote(ctx, inMint, outMint, inUnit, opts)
	if err != nil {
		return 0, err
	}
	outAtoms, err := strconv.ParseFloat(resp.OutAmount, 64)
	if err != nil {
		return 0, fmt.Errorf("parse OutAmount: %w", err)
	}
	outUnit, err := jupiter.UnitAmountByMint(outMint) // 10^decimals(OUT)
	if err != nil {
		return 0, err
	}
	return outAtoms / float64(outUnit), nil
}

// TradingFee: Jupiter комиссия 0 (только сеть).
func (j *JupiterAdapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close: дополнительных ресурсов нет.
func (j *JupiterAdapter) Close() error { return nil }
