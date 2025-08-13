package interfaces

import "github.com/gagliardetto/solana-go"

// TransactionSigner абстрагирует подпись транзакции Solana.
type TransactionSigner interface {
	PublicKey() solana.PublicKey
	SignTransaction(tx *solana.Transaction) error
}
