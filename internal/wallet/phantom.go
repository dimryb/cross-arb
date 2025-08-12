package wallet

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/gagliardetto/solana-go"
)

// PhantomWallet представляет кошелек Phantom с привязанным приватным ключом.
// Используется для управления транзакциями и ключами в сети Solana.
type PhantomWallet struct {
	privateKey solana.PrivateKey
}

var (
	ErrInvalidPrivateKey = errors.New("некорректный приватный ключ")
	ErrSignatureFailed   = errors.New("не удалось подписать")
)

// PublicKey возвращает публичный ключ кошелька.
func (w *PhantomWallet) PublicKey() solana.PublicKey {
	return w.privateKey.PublicKey()
}

// SignTransaction подписывает транзакцию Solana, используя приватный ключ кошелька.
func (w *PhantomWallet) SignTransaction(tx *solana.Transaction) error {
	if tx == nil {
		return fmt.Errorf("транзакция не может быть nil")
	}

	_, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if w.PublicKey().Equals(key) {
			return &w.privateKey
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSignatureFailed, err)
	}

	return nil
}

// NewPhantomWallet создает новый Phantom кошелек из приватного ключа.
func NewPhantomWallet(privateKeyB58 string) (*PhantomWallet, error) {
	if privateKeyB58 == "" {
		return nil, ErrInvalidPrivateKey
	}

	privateKey, err := solana.PrivateKeyFromBase58(privateKeyB58)
	if err != nil {
		return nil, err
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("некорректный размер приватного ключа: ожидается %d, получено %d",
			ed25519.PrivateKeySize, len(privateKey))
	}

	return &PhantomWallet{privateKey: privateKey}, nil
}
