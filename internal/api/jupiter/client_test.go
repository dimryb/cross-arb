package jupiter

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter/config"
	"github.com/dimryb/cross-arb/internal/logger"
)

const (
	inputMint  = "So11111111111111111111111111111111111111112"  // SOL
	outputMint = "Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB" // USDT
)

// NewTestClient создает клиент для тестов.
func NewTestClient(logger *logger.Logger) *Client {
	baseURL, _ := url.Parse(config.BASE_URL)
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

// TestClient_Quote_Integration выполняет интеграционный тест, проверяя "счастливый путь" и кастомные опции.
func TestClient_Quote_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме.")
	}

	testLogger := logger.New("debug")
	client := NewTestClient(testLogger)

	amount := int64(100000000) // 0.1 SOL

	t.Run("Default options", func(t *testing.T) {
		t.Logf("Запрашиваю котировку (опции по умолчанию) для обмена %d лампортов %s на %s...", amount, inputMint, outputMint)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		quoteResponse, err := client.Quote(ctx, inputMint, outputMint, amount, nil)

		if err != nil {
			t.Fatalf("Quote() не удалось выполнить: %v", err)
		}
		if quoteResponse.OutAmount == "" || quoteResponse.OutAmount == "0" {
			t.Errorf("Quote() вернул некорректное поле OutAmount: %s", quoteResponse.OutAmount)
		}
		if len(quoteResponse.RoutePlan) == 0 {
			t.Error("Quote() должен вернуть хотя бы один шаг в маршруте")
		}

		t.Logf("Котировка успешно получена. Сумма к получению: %s %s", quoteResponse.OutAmount, outputMint)
	})

	t.Run("Custom options", func(t *testing.T) {
		slippage := 200
		onlyDirect := true
		opts := &QuoteOptions{
			SlippageBps:      &slippage,
			OnlyDirectRoutes: &onlyDirect,
		}

		t.Logf("Запрашиваю котировку (кастомные опции) для обмена %d лампортов %s на %s...",
			amount, inputMint, outputMint)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		quoteResponse, err := client.Quote(ctx, inputMint, outputMint, amount, opts)
		if err != nil {
			t.Fatalf("Quote() с кастомными опциями не удалось выполнить: %v", err)
		}
		if quoteResponse.SlippageBps != slippage {
			t.Errorf("Ожидалось slippageBps %d, получено %d", slippage, quoteResponse.SlippageBps)
		}

		t.Logf("Котировка с кастомными опциями успешно получена. OutAmount: %s", quoteResponse.OutAmount)
	})
}

// TestClient_Quote_ValidationErrors проверяет ошибки валидации входных данных.
func TestClient_Quote_ValidationErrors(t *testing.T) {
	testLogger := logger.New("error")
	client := NewTestClient(testLogger)
	ctx := context.Background()

	testCases := []struct {
		name       string
		inputMint  string
		outputMint string
		amount     int64
		wantErrMsg string
	}{
		{"Empty inputMint", "", outputMint, 1, "inputMint cannot be empty"},
		{"Empty outputMint", inputMint, "", 1, "outputMint cannot be empty"},
		{"Zero amount", inputMint, outputMint, 0, "amount must be positive"},
		{"Negative amount", inputMint, outputMint, -1, "amount must be positive"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := client.Quote(ctx, tc.inputMint, tc.outputMint, tc.amount, nil)
			if err == nil || !strings.Contains(err.Error(), tc.wantErrMsg) {
				t.Errorf("Ожидалась ошибка, содержащая '%s', но получено: %v", tc.wantErrMsg, err)
			}
		})
	}
}

// TestClient_Quote_ContextTimeout проверяет обработку отмены контекста.
func TestClient_Quote_ContextTimeout(t *testing.T) {
	testLogger := logger.New("error")
	client := NewTestClient(testLogger)

	// Создаем контекст, который отменится через 1 наносекунду
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	_, err := client.Quote(ctx, inputMint, outputMint, 1000, nil)

	if err == nil {
		t.Fatal("Ожидалась ошибка отмены контекста, но ее не было")
	}

	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Ожидалась ошибка 'context deadline exceeded', получено: %v", err)
	}

	t.Logf("Тест на таймаут контекста успешно пройден с ошибкой: %v", err)
}
