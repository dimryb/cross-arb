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
const tokenListURL = "https://tokens.jup.ag/tokens?tags=verified,community,strict" // #nosec G101

type TokenEntry struct {
	Symbol  string `json:"symbol"`
	Address string `json:"address"`
	// Decimals число знаков после запятой у котируемого токена (USDT = 6).
	Decimals uint8 `json:"decimals"`
}

var (
	tokenOnce      sync.Once
	tokenMap       map[string]TokenEntry // SYMBOL -> entry
	tokenMapByMint map[string]TokenEntry // MINT ADDRESS -> entry
	tokenErr       error
)

// UnitAmount возвращает int64, равный 1*10^decimals для заданного тикера.
func UnitAmount(ticker string) (int64, error) {
	dec, err := getDecimals(ticker)
	if err != nil {
		return 0, err
	}
	return int64(math.Pow10(int(dec))), nil
}

// UnitAmount возвращает int64, равный 1*10^decimals для заданного тикера.
func UnitAmountByMint(mint string) (int64, error) {
	dec, err := getDecimalsByMint(mint)
	if err != nil {
		return 0, err
	}
	return int64(math.Pow10(int(dec))), nil
}

// getDecimalsByMint возвращает decimals по mint-адресу (как в списке Jupiter).
func getDecimalsByMint(mint string) (uint8, error) {
	m, err := getTokenMapByMint()
	if err != nil {
		return 0, err
	}
	key := strings.TrimSpace(mint) // регистр в base58 важен — не меняем
	entry, ok := m[key]
	if !ok {
		return 0, fmt.Errorf("decimals not found for mint %s", mint)
	}
	return entry.Decimals, nil
}

// getDecimals возвращает количество знаков после запятой (decimals) токена по символу.
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

// getTokenMap загружает (один раз) список токенов Jupiter и строит map[SYMBOL]TokenEntry
// и map[MINT]TokenEntry.
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

		bySymbol := make(map[string]TokenEntry, len(list))
		byMint := make(map[string]TokenEntry, len(list))
		for _, t := range list {
			bySymbol[strings.ToUpper(t.Symbol)] = t
			byMint[strings.TrimSpace(t.Address)] = t // адрес оставляем в исходном регистре
		}
		tokenMap = bySymbol
		tokenMapByMint = byMint
	})
	return tokenMap, tokenErr
}

// getTokenMapByMint возвращает карту mint->entry, гарантируя, что список загружен.
func getTokenMapByMint() (map[string]TokenEntry, error) {
	if tokenMapByMint != nil || tokenErr != nil {
		return tokenMapByMint, tokenErr
	}
	// Инициализируем через общий загрузчик
	_, err := getTokenMap()
	if err != nil {
		return nil, err
	}
	return tokenMapByMint, tokenErr
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
var quoteCurrencies = []string{
	"USDT",
	"USDC",
	"BUSD",
	"USD",
}
