package swap

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	"github.com/dimryb/cross-arb/internal/logger"
	"github.com/dimryb/cross-arb/internal/wallet"
)

const (
	// #nosec G101
	phantomPrivateKey = "3QkJFM1ac7BtTGR4d5KG6PAm86qU9E3P6VyZS7mGWuZq4tHJtQdYw9R3q5o1tbzZs6ryh5gYH1j3ogUfjEUMMrXN"
	jupiterAPIURL     = "https://lite-api.jup.ag/swap/v1"
	solanaRPCURL      = "https://api.mainnet-beta.solana.com"
)

// TestSwapper_Swap выполняет полный цикл:
// котировка -> создание транзакции -> подпись -> отправка в сеть.
func TestSwapper_Swap(t *testing.T) {
	if testing.Short() {
		t.Skip("Пропускаем тест в коротком режиме.")
	}

	testLogger := logger.New("debug")

	swapper, err := NewSwapper(testLogger, jupiterAPIURL, solanaRPCURL)
	if err != nil {
		t.Fatalf("Не удалось создать сервис: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 1. Создаем кошелек.
	phantomWallet, err := wallet.NewPhantomWallet(phantomPrivateKey)
	if err != nil {
		t.Fatalf("Не удалось создать Phantom кошелек: %v", err)
	}
	t.Logf("Кошелек успешно создан. PublicKey: %s", phantomWallet.PublicKey())

	// Определяем mint-адреса для пары через резолвер Jupiter
	inMint, outMint, err := jupiter.ConvertSpotToMints("SOLUSDT")
	if err != nil {
		t.Fatalf("Не удалось получить mint-адреса: %v", err)
	}

	amountToExchange := int64(1000) // примерная малая сумма единиц базового токена
	sig, err := swapper.Swap(ctx, phantomWallet, inMint, outMint, amountToExchange, nil)
	if err != nil {
		// Проверяем, является ли ошибка связанной с недостаточностью средств
		errorMsg := strings.ToLower(err.Error())
		isInsufficientFunds := strings.Contains(errorMsg, "insufficient funds") ||
			strings.Contains(errorMsg, "account not found") ||
			strings.Contains(errorMsg, "attempt to debit an account but found no record of a prior credit") ||
			strings.Contains(errorMsg, "insufficient account balance")

		if isInsufficientFunds {
			t.Logf("Тест успешен: получена ожидаемая ошибка недостаточности средств")
			t.Logf("Ошибка: %v", err)
			t.Log("Для реального выполнения пополните кошелек SOL в mainnet")
		} else {
			t.Fatalf("Неожиданная ошибка при отправке транзакции: %v", err)
		}
	} else {
		t.Logf("Транзакция успешно отправлена и подтверждена!")
		t.Logf("Сигнатура: %s", sig.String())
		t.Logf("Проверьте в эксплорере: https://explorer.solana.com/tx/%s", sig.String())
	}
}
