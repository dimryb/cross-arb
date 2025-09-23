package jupiter

import (
	"context"
	"fmt"
	"math"
	"math/big"
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

// toAtoms переводит человеко-понятный объём в атомы токена с округлением к ближайшему целому атому.
// toAtomsExactUint переводит человеко-понятный объём в атомы токена с округлением half-up и возвращает uint64.
func toAtomsExactUint(amount float64, unit uint64) (uint64, error) {
	if math.IsNaN(amount) || math.IsInf(amount, 0) {
		return 0, fmt.Errorf("invalid amount: %v", amount)
	}
	r := new(big.Rat).SetFloat64(amount)
	if r == nil {
		return 0, fmt.Errorf("cannot represent amount as rational: %v", amount)
	}
	r.Mul(r, new(big.Rat).SetInt(new(big.Int).SetUint64(unit)))

	num := new(big.Int).Set(r.Num())
	den := new(big.Int).Set(r.Denom())
	if den.Sign() == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	neg := num.Sign() < 0
	if neg {
		num.Neg(num)
	}
	q, rem := new(big.Int).QuoRem(num, den, new(big.Int))
	twiceRem := new(big.Int).Lsh(rem, 1)
	if twiceRem.Cmp(den) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if neg {
		q.Neg(q)
	}
	if q.Sign() < 0 || !q.IsUint64() {
		return 0, fmt.Errorf("atom amount overflow or negative")
	}
	return q.Uint64(), nil
}

// TradingFee Jupiter комиссия 0 (только сеть).
func (j *Adapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close дополнительных ресурсов нет.
func (j *Adapter) Close() error { return nil }
