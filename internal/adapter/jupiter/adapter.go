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

// Name удовлетворяет интерфейсу EXAdapter.
func (j *Adapter) Name() string { return "jupiter" }

// Quote для Jupiter: цены в QUOTE за 1 BASE для пары, например SOL/USDT.
func (j *Adapter) Quote(
	ctx context.Context,
	pair string,
	baseAmount int64,
) (float64, float64, error) {
	mints, ok := j.pairConfig[pair]
	if !ok {
		return 0, 0, fmt.Errorf("неизвестная пара %s", pair)
	}

	resp, err := j.client.Quote(ctx, mints.InputMint, mints.OutputMint, baseAmount, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при получении данных api: %w", err)
	}

	inAmount, err := strconv.ParseFloat(resp.InAmount, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования параметра inAmount к float64: %w", err)
	}
	outAmount, err := strconv.ParseFloat(resp.OutAmount, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка преобразования параметра outAmount к float64: %w", err)
	}

	ask := inAmount / outAmount
	bid := outAmount / inAmount

	return ask, bid, nil
}

// TradingFee Jupiter комиссия 0 (только сеть).
func (j *Adapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close дополнительных ресурсов нет.
func (j *Adapter) Close() error { return nil }
