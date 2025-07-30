package jupiter

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/dimryb/cross-arb/internal/logger"
)

const (
	userPublicKey = "EQQp359jRUSmzTSdrHVzKK7mTfTLompX637j9yLahQkb" // Devnet-кошелек
)

// Константы для теста TestClient_Swap_WithAllParameters_Integration.
const (
	// Используем один тестовый адрес для всех опциональных полей аккаунтов.
	optionalAccount = "7g166K8tY855H43GN2mR5z3f61xTfVzWdE6aT4xR3JgY"
	// Приоритетная комиссия в микро-лампортах (5e-11 SOL).
	computeUnitPriceMicroLamports uint64 = 50000
	// Срок жизни blockhash в слотах.
	blockhashSlotsToExpiry = 100
)

// TestClient_Swap_Integration выполняет базовый интеграционный тест для /swap.
func TestClient_Swap_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме.")
	}

	testLogger := logger.New("debug")
	client := NewTestClient(testLogger)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Получаем актуальную котировку
	quoteResponse, err := client.Quote(ctx, inputMint, outputMint, amountToExchange, nil)
	if err != nil {
		t.Fatalf("Не удалось получить котировку для теста обмена: %v", err)
	}

	// 2. Создаем запрос на обмен
	swapReq := &SwapRequest{
		QuoteResponse: *quoteResponse,
		UserPublicKey: userPublicKey,
	}

	// 3. Выполняем обмен
	swapResponse, err := client.Swap(ctx, swapReq)
	if err != nil {
		t.Fatalf("Swap() не удалось выполнить: %v", err)
	}

	// 4. Проверяем результат
	if swapResponse.SwapTransaction == "" {
		t.Error("Swap() вернул пустую транзакцию")
	}

	t.Logf("Транзакция для обмена успешно создана. LastValidBlockHeight: %d", swapResponse.LastValidBlockHeight)
}

// TestClient_Swap_WithAllParameters_Integration проверяет создание обмена со всеми возможными параметрами.
func TestClient_Swap_WithAllParameters_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме.")
	}
	testLogger := logger.New("debug")
	client := NewTestClient(testLogger)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	// 1. Получаем котировку с поддержкой legacy transaction
	quoteOpts := &QuoteOptions{AsLegacyTransaction: boolPtr(true)}
	quoteResponse, err := client.Quote(ctx, inputMint, outputMint, amountToExchange, quoteOpts)
	if err != nil {
		t.Fatalf("Не удалось получить котировку для теста: %v", err)
	}

	// 2. Создаем запрос на обмен со всеми параметрами, используя константы.
	swapReq := &SwapRequest{
		// --- Обязательные поля ---
		QuoteResponse: *quoteResponse,
		UserPublicKey: userPublicKey,
		// --- Опциональные поля ---
		Payer:                         stringPtr(userPublicKey), // Плательщик комиссий, обычно совпадает с UserPublicKey
		WrapAndUnwrapSol:              boolPtr(true),
		UseSharedAccounts:             boolPtr(true),
		FeeAccount:                    stringPtr(optionalAccount), // Пример аккаунта для сбора комиссии платформы
		TrackingAccount:               stringPtr(optionalAccount), // Пример аккаунта для трекинга
		DestinationTokenAccount:       stringPtr(optionalAccount), // Пример кастомного аккаунта для получения токенов
		AsLegacyTransaction:           boolPtr(true),              // Должно совпадать с запросом котировки
		DynamicComputeUnitLimit:       boolPtr(true),
		SkipUserAccountsRPCCalls:      boolPtr(true),
		DynamicSlippage:               boolPtr(true),
		ComputeUnitPriceMicroLamports: uint64Ptr(computeUnitPriceMicroLamports), // Приоритетная комиссия
		BlockhashSlotsToExpiry:        intPtr(blockhashSlotsToExpiry),           // Срок жизни blockhash
		// `PrioritizationFeeLamports` не задается, так как `ComputeUnitPriceMicroLamports` уже задан
	}

	// 3. Выполняем обмен
	t.Log("Отправка запроса на обмен со всеми параметрами...")
	swapResponse, err := client.Swap(ctx, swapReq)
	if err != nil {
		t.Fatalf("Swap() со всеми параметрами не удалось выполнить: %v", err)
	}

	// 4. Проверяем результат, если ошибки не было
	if swapResponse.SwapTransaction == "" {
		t.Error("Swap() со всеми параметрами вернул пустую транзакцию")
	}

	t.Logf("Транзакция со всеми параметрами успешно создана. LastValidBlockHeight: %d", swapResponse.LastValidBlockHeight)
}

// TestClient_Swap_ValidationErrors проверяет ошибки валидации для метода Swap.
func TestClient_Swap_ValidationErrors(t *testing.T) {
	client := NewTestClient(logger.New("error"))
	ctx := context.Background()

	swapReq := &SwapRequest{UserPublicKey: ""}

	_, err := client.Swap(ctx, swapReq)
	if err == nil || !strings.Contains(err.Error(), "UserPublicKey не может быть пустым") {
		t.Errorf("Ожидалась ошибка валидации для UserPublicKey, но получено: %v", err)
	}
}

// TestClient_Swap_ContextTimeout проверяет обработку отмены контекста.
func TestClient_Swap_ContextTimeout(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем интеграционный тест в коротком режиме.")
	}
	client := NewTestClient(logger.New("error"))

	quoteResponse, err := client.Quote(context.Background(), inputMint, outputMint, amountToExchange, nil)
	if err != nil {
		t.Fatalf("Не удалось получить котировку для теста: %v", err)
	}

	// Создаем контекст, который отменится немедленно
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	time.Sleep(1 * time.Millisecond) // Даем контексту время на отмену

	swapReq := &SwapRequest{
		QuoteResponse: *quoteResponse,
		UserPublicKey: userPublicKey,
	}

	_, err = client.Swap(ctx, swapReq)
	if err == nil || !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("Ожидалась ошибка 'context deadline exceeded', получено: %v", err)
	}
}

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uint64Ptr(u uint64) *uint64 { return &u }
func intPtr(i int) *int          { return &i }
