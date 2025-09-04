package jupiter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	i "github.com/dimryb/cross-arb/internal/interface"
)

// Важно: тип MintPair объявлен в internal/adapter/jupiter/mint_pair.go
// type MintPair struct { BaseMint, QuoteMint string }

// AdapterConfig конфигурация адаптера.
type AdapterConfig struct {
	BaseURL string
	Enabled bool
	Timeout time.Duration
	Pairs   map[string]MintPair // необязательная карта маппинга "SOL/USDT" → мины
}

// Adapter реализует DEXAdapter через Jupiter.
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
	cli, err := jupiter.NewJupiterClient(l, cfg.BaseURL)
	if err != nil {
		l.Fatalf("failed to create Jupiter client: %v", err)
	}

	var pairs map[string]MintPair
	if cfg != nil && cfg.Pairs != nil {
		pairs = cfg.Pairs
	} else {
		pairs = make(map[string]MintPair)
	}

	return &Adapter{
		client:     cli,
		logger:     l.Named("jupiter"),
		baseURL:    cfg.BaseURL,
		pairConfig: pairs,
	}
}

// Name удовлетворяет интерфейсу EXAdapter.
func (j *Adapter) Name() string { return "jupiter" }

// Quote возвращает эффективные котировки bid/ask (QUOTE за 1 BASE) для заданного объёма baseAmount (в BASE).
// bid — сколько QUOTE вы получите, если ПРОДАДИТЕ baseAmount BASE (ExactIn: BASE→QUOTE).
// ask — сколько QUOTE вам потребуется, чтобы КУПИТЬ baseAmount BASE (аппроксимация через поиск: QUOTE→BASE).
func (j *Adapter) Quote(
	ctx context.Context,
	pair string,
	baseAmount float64,
) (bid, ask float64, err error) {
	if ctx == nil {
		return 0, 0, fmt.Errorf("nil context")
	}
	if baseAmount <= 0 || math.IsNaN(baseAmount) || math.IsInf(baseAmount, 0) {
		return 0, 0, fmt.Errorf("invalid baseAmount: %v", baseAmount)
	}
	mints, err := j.resolveMints(pair)
	if err != nil {
		return 0, 0, err
	}

	// Сначала считаем bid: ExactIn BASE->QUOTE на нужный объём.
	inUnit, err := jupiter.UnitAmountByMint(mints.BaseMint) // 10^decimals(BASE)
	if err != nil {
		return 0, 0, err
	}
	inAtoms := int64(math.Round(baseAmount * float64(inUnit)))
	if inAtoms <= 0 {
		return 0, 0, fmt.Errorf("amount too small after conversion to atoms")
	}

	respBid, err := j.client.Quote(ctx, mints.BaseMint, mints.QuoteMint, inAtoms, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("jupiter quote (BASE->QUOTE): %w", err)
	}
	if respBid == nil {
		return 0, 0, errors.New("empty quote response (BASE->QUOTE)")
	}

	outAtomsBid, err := strconv.ParseFloat(respBid.OutAmount, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("parse OutAmount for bid: %w", err)
	}
	outUnitQuote, err := jupiter.UnitAmountByMint(mints.QuoteMint) // 10^decimals(QUOTE)
	if err != nil {
		return 0, 0, err
	}
	quoteOut := outAtomsBid / float64(outUnitQuote)
	bid = quoteOut / baseAmount // QUOTE per 1 BASE

	// Далее считаем ask: сколько QUOTE нужно, чтобы получить baseAmount BASE (ExactIn QUOTE->BASE c поиском).
	ask, err = j.estimateAsk(ctx, mints, baseAmount, bid)
	if err != nil {
		return 0, 0, fmt.Errorf("estimate ask: %w", err)
	}
	return bid, ask, nil
}

// estimateAsk оценивает ask (QUOTE за 1 BASE), решая задачу:
// найти inQuote (в QUOTE), чтобы QUOTE->BASE (ExactIn) дал outBase ~= baseAmount.
// Поиск — экспоненциальный рост верхней границы + бинарный поиск.
func (j *Adapter) estimateAsk(
	ctx context.Context,
	mints MintPair,
	baseAmount float64,
	bidHint float64,
) (float64, error) {
	// Границы поиска (в QUOTE). Начальная оценка — hint от bid.
	lower := 0.0
	upper := math.Max(bidHint*baseAmount*1.10, 1e-9)

	// f(x): сколько BASE получим, если отдадим x QUOTE.
	f := func(x float64) (float64, error) {
		inUnitQ, err := jupiter.UnitAmountByMint(mints.QuoteMint)
		if err != nil {
			return 0, err
		}
		inAtoms := int64(math.Round(x * float64(inUnitQ)))
		if inAtoms <= 0 {
			return 0, fmt.Errorf("amount too small after conversion to atoms")
		}
		resp, err := j.client.Quote(ctx, mints.QuoteMint, mints.BaseMint, inAtoms, nil)
		if err != nil {
			return 0, err
		}
		if resp == nil {
			return 0, fmt.Errorf("empty quote response")
		}
		outAtoms, err := strconv.ParseFloat(resp.OutAmount, 64)
		if err != nil {
			return 0, fmt.Errorf("parse OutAmount: %w", err)
		}
		outUnitB, err := jupiter.UnitAmountByMint(mints.BaseMint)
		if err != nil {
			return 0, err
		}
		return outAtoms / float64(outUnitB), nil
	}

	// Увеличиваем верхнюю границу пока outBase < baseAmount.
	var out float64
	for i := 0; i < 8; i++ {
		var err error
		out, err = f(upper)
		if err != nil {
			return 0, err
		}
		if out >= baseAmount {
			break
		}
		upper *= 2
	}
	if out < baseAmount {
		if upper > 0 {
			return upper / baseAmount, nil
		}
		return 0, fmt.Errorf("unable to bound ask for %f BASE", baseAmount)
	}

	// Бинарный поиск.
	for i := 0; i < 16; i++ {
		mid := (lower + upper) / 2
		outMid, err := f(mid)
		if err != nil {
			return 0, err
		}
		if math.Abs(outMid-baseAmount) <= baseAmount*1e-6 {
			return mid / baseAmount, nil
		}
		if outMid < baseAmount {
			lower = mid
		} else {
			upper = mid
		}
	}
	return upper / baseAmount, nil
}

// quote возвращает: «сколько OUT токенов за 1 IN токен» в единицах OUT (не в атомах).
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
	// 1 IN в атомах.
	inAtoms := inUnit

	resp, err := j.client.Quote(ctx, inMint, outMint, inAtoms, opts)
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

// resolveMints: если пара есть в конфиге — берём её; иначе "BASE/QUOTE" → "BASEQUOTE" → mints через токен-резолвер.
func (j *Adapter) resolveMints(pair string) (MintPair, error) {
	if m, ok := j.pairConfig[pair]; ok && m.BaseMint != "" && m.QuoteMint != "" {
		return m, nil
	}
	s := strings.ToUpper(strings.TrimSpace(pair))
	s = strings.ReplaceAll(s, "/", "")
	in, out, err := jupiter.ConvertSpotToMints(s)
	if err != nil {
		return MintPair{}, fmt.Errorf("resolve mints for pair %s: %w", pair, err)
	}
	return MintPair{BaseMint: in, QuoteMint: out}, nil
}

// TradingFee Jupiter комиссия 0 (только сеть).
func (j *Adapter) TradingFee(string) (maker, taker float64) { return 0, 0 }

// Close дополнительных ресурсов нет.
func (j *Adapter) Close() error { return nil }
