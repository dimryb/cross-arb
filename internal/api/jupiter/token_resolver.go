package jupiter

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)

// tokenListURL эндпоинт для запроса токенов для Jupiter.
const tokenListURL = "https://tokens.jup.ag/tokens?tags=verified,community,strict"

type TokenEntry struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`

	// Decimals число знаков после запятой у котируемого токена (USDT = 6).
	Decimals uint8 `json:"decimals"`
}

var (
	tokenOnce sync.Once
	tokenMap  map[string]TokenEntry
	tokenErr  error
)

// UnitAmount возвращает int64, равный 1*10^decimals для заданного тикера.
func UnitAmount(ticker string) (int64, error) {
	dec, err := getDecimals(ticker)
	if err != nil {
		return 0, err
	}
	return int64(math.Pow10(int(dec))), nil
}

// getDecimals возвращает количество знаков после запятой (decimals) токена.
func getDecimals(ticker string) (uint8, error) {
	m, err := getTokenMap()
	if err != nil {
		return 0, err
	}

	entry, ok := m[strings.ToUpper(strings.TrimSpace(ticker))]
	if !ok {
		return 0, fmt.Errorf("decimals not found for token %s", ticker)
	}
	return entry.Decimals, nil
}

// getTokenMap загружает (один раз) список токенов Jupiter и строит map[SYMBOL]TokenEntry.
func getTokenMap() (map[string]TokenEntry, error) {
	tokenOnce.Do(func() {
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tokenListURL, nil)
		if err != nil {
			tokenErr = fmt.Errorf("build request: %w", err)
			return
		}

		resp, err := client.Do(req)
		if err != nil {
			tokenErr = fmt.Errorf("fetch token list: %w", err)
			return
		}

		if resp.StatusCode != http.StatusOK {
			tokenErr = fmt.Errorf("unexpected status %d when fetching token list", resp.StatusCode)
			_ = resp.Body.Close()
			return
		}
		defer func() { _ = resp.Body.Close() }()

		var list []TokenEntry
		if err = json.NewDecoder(resp.Body).Decode(&list); err != nil {
			tokenErr = fmt.Errorf("decode token list: %w", err)
			return
		}

		tmp := make(map[string]TokenEntry, len(list))
		for _, t := range list {
			tmp[strings.ToUpper(t.Symbol)] = t
		}
		tokenMap = tmp
	})
	return tokenMap, tokenErr
}

// ConvertSpotToMints преобразует символ CEX-биржи вида "SOLUSDT" в inputMint/outputMint.
func ConvertSpotToMints(symbol string) (string, string, error) {
	base, quote, err := ParseSpotSymbol(symbol)
	if err != nil {
		return "", "", err
	}

	in, err := getMint(base)
	if err != nil {
		return "", "", err
	}

	out, err := getMint(quote)
	if err != nil {
		return "", "", err
	}

	return in, out, nil
}

// ParseSpotSymbol разбивает CEX-символ формата "SOLUSDT" на базовый и котировочный токены.
func ParseSpotSymbol(symbol string) (base, quote string, err error) {
	s := strings.ToUpper(strings.TrimSpace(symbol))
	for _, q := range quoteCurrencies {
		if strings.HasSuffix(s, q) {
			return strings.TrimSuffix(s, q), q, nil
		}
	}
	return "", "", fmt.Errorf("unrecognized spot symbol: %s", symbol)
}

// getMint возвращает mint-адрес по тикеру (например "USDT").
func getMint(ticker string) (string, error) {
	m, err := getTokenMap()
	if err != nil {
		return "", err
	}

	entry, ok := m[strings.ToUpper(strings.TrimSpace(ticker))]
	if !ok {
		return "", fmt.Errorf("mint not found for token %s", ticker)
	}
	return entry.Address, nil
}

// quoteCurrencies содержит наиболее популярные суффиксы котировочных валют.
// Упорядочены по длине, чтобы "USDT" не бился как "USD".
var quoteCurrencies = []string{
	"USDT",
	"USDC",
	"BUSD",
	"USD",
}
