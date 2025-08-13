package swap

import (
	"context"
	"errors"
	"fmt"

	"github.com/dimryb/cross-arb/internal/api/jupiter"
	i "github.com/dimryb/cross-arb/internal/interface"
	"github.com/gagliardetto/solana-go"
)

// SwapWithQuote выполняет обмен, используя уже полученную котировку Jupiter.
//
// Ожидается, что котировка (quote) ещё валидна на момент вызова. Метод:
//  1. запрашивает у Jupiter сериализованную транзакцию для обмена,
//  2. валидирует соответствие fee payer ожидаемому подписанту и наличие blockhash,
//  3. подписывает транзакцию переданным TransactionSigner,
//  4. отправляет и подтверждает её через Solana RPC.
//
// Возвращает сигнатуру подтверждённой транзакции или ошибку.
func (s *Swapper) SwapWithQuote(
	ctx context.Context,
	signer i.TransactionSigner,
	quote *jupiter.QuoteResponse,
) (solana.Signature, error) {
	swapReq := &jupiter.SwapRequest{
		QuoteResponse: *quote,
		UserPublicKey: signer.PublicKey().String(),
	}

	swapResp, err := s.apiClient.Swap(ctx, swapReq)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("не удалось создать транзакцию: %w", err)
	}

	tx, err := solana.TransactionFromBase64(swapResp.SwapTransaction)
	if err != nil {
		return solana.Signature{}, fmt.Errorf("не удалось десериализовать транзакцию из base64: %w", err)
	}

	if err := validateTransactionForSigner(tx, signer.PublicKey()); err != nil {
		return solana.Signature{}, fmt.Errorf("валидация транзакции перед подписью не пройдена: %w", err)
	}

	if err := signer.SignTransaction(tx); err != nil {
		return solana.Signature{}, fmt.Errorf("не удалось подписать транзакцию: %w", err)
	}

	signature, err := s.solanaClient.SendAndConfirmTransaction(ctx, tx)
	if err != nil {
		return solana.Signature{}, err
	}

	return signature, nil
}

// validateTransactionForSigner проверяет, что транзакция готова к подписи данным сайнером.
// Валидирует непустой message/blockhash и совпадение fee payer с ожидаемым публичным ключом.
func validateTransactionForSigner(tx *solana.Transaction, expectedSigner solana.PublicKey) error {
	if tx == nil {
		return errors.New("tx is nil")
	}
	// Проверяем наличие recent blockhash (legacy) или валидного message (v0).
	if tx.Message.RecentBlockhash.IsZero() && len(tx.Message.AccountKeys) == 0 {
		return errors.New("transaction has empty message or blockhash")
	}
	// Fee payer — первый аккаунт в списке подписантов.
	if len(tx.Message.AccountKeys) == 0 {
		return errors.New("transaction has no account keys")
	}
	feePayer := tx.Message.AccountKeys[0]
	if !feePayer.Equals(expectedSigner) {
		return fmt.Errorf("fee payer %s не совпадает с ожидаемым подписантом %s", feePayer.String(), expectedSigner.String())
	}
	return nil
}
