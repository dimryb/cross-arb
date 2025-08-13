package swap

import (
	"context"
	"fmt"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/gagliardetto/solana-go"
)

// Swap — удобная обёртка над полным циклом обмена без предварительно переданной котировки.
// Метод сам запрашивает котировку у Jupiter, после чего делегирует выполнение в SwapWithQuote.
//
// Возвращает сигнатуру подтверждённой транзакции или ошибку.
func (s *Swapper) Swap(
	ctx context.Context,
	signer i.TransactionSigner,
	inputMint, outputMint string,
	amount int64,
	opts *jupiter.QuoteOptions,
) (solana.Signature, error) {
	quote, err := s.apiClient.Quote(ctx, inputMint, outputMint, amount, opts)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("не удалось получить котировку: %w", err)
	}
	return s.SwapWithQuote(ctx, signer, quote)
}
