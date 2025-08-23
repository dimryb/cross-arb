package jupiter

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// AdapterConfig конфигурация адаптера.
type AdapterConfig struct {
	BaseURL string
	Enabled bool
	Timeout time.Duration       // Может использоваться в NewJupiterAdapterFromConfig
	Pairs   map[string]MintPair // symbol → [base_mint, quote_mint]
}

// Adapter JupiterAdapter использует публичный Quote-API агрегатора Jupiter (Solana).
// Он запрашивает цену обмена base → quote (ask) и quote → base (bid).
type Adapter struct {
	client     *jupiter.Client
	logger     i.Logger
	baseURL    string
	pairConfig map[string]MintPair // "SOL/USDT" → {baseMint, quoteMint}
}

// NewAdapter создаёт адаптер.
// pairMap: "SOL/USDT": {baseMint, quoteMint}.
func NewAdapter(l i.Logger, cfg *AdapterConfig) *Adapter {
	// TODO: handle error
	client, err := jupiter.NewJupiterClient(l, cfg.BaseURL)
	if err != nil {
		// Линтер errcheck: ошибку не игнорируем
		l.Fatalf("failed to create Jupiter client: %v", err)
	}

	return &Adapter{
		client:     client,
		logger:     l.Named("jupiter"),
		baseURL:    cfg.BaseURL,
		pairConfig: cfg.Pairs,
	}
}

// Name удовлетворяет интерфейсу ExchangeAdapter.
func (j *Adapter) Name() string { return "jupiter" }

// OrderBookTop для Jupiter: цены в QUOTE за 1 BASE для пары, например SOL/USDT.
func (j *Adapter) OrderBookTop(
	ctx context.Context,
	pair string,
) (bestBid, bestAsk float64, err error) {
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
func (j *Adapter) quote(
	ctx context.Context,
	inMint string,
	outMint string,
	opts *jupiter.QuoteOptions,
) (float64, error) {
	inUnit, err := jupiter.UnitAmountByMint(inMint) // 10^decimals(IN)
	if err != nil {
		return 0, err
	}

	resp, err := j.client.Quote(ctx, inMint, outMint, inUnit, opts)
	if err != nil {
		return 0, err
	}
	if resp == nil {
		return 0, fmt.Errorf("empty response from jupiter")
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

// TradingFee Jupiter комиссия 0 (только сеть).
func (j *Adapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close дополнительных ресурсов нет.
func (j *Adapter) Close() error { return nil }
