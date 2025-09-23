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

// Quote возвращает эффективные котировки bid/ask (QUOTE per BASE) для указанного объёма BASE
// в человеко-понятных единицах (например, 1.0 SOL). Внутри метод конвертирует объём BASE
// в атомы и вызывает Jupiter Client, который принимает атомы:
//   - bid: ExactIn BASE->QUOTE
//   - ask: ExactOut QUOTE->BASE
func (j *Adapter) Quote(
	ctx context.Context,
	pair string,
	baseAmount float64,
) (float64, float64, error) {
	mints, ok := j.pairConfig[pair]
	if !ok {
		return 0, 0, fmt.Errorf("неизвестная пара %s", pair)
	}

	// 1) Узнаём размерность токенов (сколько атомов в 1 BASE/QUOTE): 10^decimals
	baseUnits, err := jupiter.UnitAmountByMint(mints.InputMint)
	if err != nil {
		return 0, 0, fmt.Errorf("decimals base: %w", err)
	}
	quoteUnits, err := jupiter.UnitAmountByMint(mints.OutputMint)
	if err != nil {
		return 0, 0, fmt.Errorf("decimals quote: %w", err)
	}

	if baseUnits < 0 {
		return 0, 0, fmt.Errorf("отрицательный decimals: baseUnits=%d", baseUnits)
	}
	if quoteUnits < 0 {
		return 0, 0, fmt.Errorf("отрицательный decimals: quoteUnits=%d", quoteUnits)
	}

	baseUnit := uint64(baseUnits)
	quoteUnit := uint64(quoteUnits)

	// Переводим запрошенный объём BASE (человеко-читаемый объем) в атомы BASE для вызовов клиента
	baseAmountAtoms, err := toAtomsExactUint(baseAmount, baseUnit)
	if err != nil || baseAmountAtoms == 0 {
		return 0, 0, fmt.Errorf("некорректный объём baseAmount=%v: %w", baseAmount, err)
	}

	// Считаем BID через отдельный метод (ExactIn BASE atoms)
	bid, err := j.computeBid(ctx, mints, baseAmountAtoms, baseUnit, quoteUnit)
	if err != nil {
		return 0, 0, err
	}

	// 3) ASK: ExactOut QUOTE->BASE на baseAmountAtoms атомов BASE
	ask, err := j.computeAsk(ctx, mints, baseAmountAtoms, baseUnit, quoteUnit)
	if err != nil {
		return 0, 0, err
	}

	return bid, ask, nil
}

func (j *Adapter) computeAsk(
	ctx context.Context,
	mints MintPair,
	baseAmountAtoms uint64,
	baseUnit, quoteUnit uint64,
) (float64, error) {
	// computeAsk считает эффективную цену ASK (QUOTE per BASE) в режиме ExactOut,
	// рассчитывая сколько QUOTE требуется для получения фиксированного количества BASE (в атомах).
	mode := jupiter.SwapModeExactOut
	opts := jupiter.QuoteOptions{SwapMode: &mode}

	resp, err := j.client.Quote(ctx, mints.OutputMint, mints.InputMint, baseAmountAtoms, &opts)
	if err == nil {
		inQuote, err := parseU64(resp.InAmount, "inAmount(ask)")
		if err != nil {
			return 0, err
		}
		outBase, err := parseU64(resp.OutAmount, "outAmount(ask)")
		if err != nil {
			return 0, err
		}
		if outBase == 0 || inQuote == 0 {
			return 0, fmt.Errorf("пустой маршрут для ask: in=%d out=%d", inQuote, outBase)
		}
		return ratPrice(inQuote, quoteUnit, outBase, baseUnit), nil // QUOTE per BASE
	}
	return 0, fmt.Errorf("ошибка при получении данных api: %w", err)
}

// computeBid считает эффективную цену BID (QUOTE per BASE) в режиме ExactIn,
// рассчитывая сколько QUOTE получится при обмене фиксированного количества BASE (в атомах).
func (j *Adapter) computeBid(
	ctx context.Context,
	mints MintPair,
	baseAmountAtoms uint64,
	baseUnit, quoteUnit uint64,
) (float64, error) {
	mode := jupiter.SwapModeExactIn
	opts := jupiter.QuoteOptions{SwapMode: &mode}

	resp, err := j.client.Quote(ctx, mints.InputMint, mints.OutputMint, baseAmountAtoms, &opts)
	if err == nil {
		inBase, err := parseU64(resp.InAmount, "inAmount(bid)")
		if err != nil {
			return 0, err
		}
		outQuote, err := parseU64(resp.OutAmount, "outAmount(bid)")
		if err != nil {
			return 0, err
		}
		if inBase == 0 || outQuote == 0 {
			return 0, fmt.Errorf("пустой маршрут для bid: in=%d out=%d", inBase, outQuote)
		}
		// QUOTE per BASE
		return ratPrice(outQuote, quoteUnit, inBase, baseUnit), nil
	}
	return 0, fmt.Errorf("ошибка при получении данных api: %w", err)
}

// (numAtoms/numUnit) / (denAtoms/denUnit).
func ratPrice(numAtoms, numUnit, denAtoms, denUnit uint64) float64 {
	n := new(big.Int).Mul(new(big.Int).SetUint64(numAtoms), new(big.Int).SetUint64(denUnit))
	d := new(big.Int).Mul(new(big.Int).SetUint64(denAtoms), new(big.Int).SetUint64(numUnit))
	if d.Sign() == 0 {
		return math.NaN()
	}
	r := new(big.Rat).SetFrac(n, d)
	f, _ := r.Float64()
	return f
}

func parseU64(s, what string) (uint64, error) {
	u, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse %s: %w", what, err)
	}
	return u, nil
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
